package communicate

import (
	"bonbon/database"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type user struct {
	conns     []*websocket.Conn
	lock      *sync.Mutex // lock住conns
	match     int         // -1表示無配對
	matchType string
	bonbon    bool
}

var onlineUser = make(map[int]*user)
var onlineLock = new(sync.RWMutex)

// 粒度高，將降低效能
// TODO: 改為讀寫鎖或拆分粒度
var globalMatchLock = new(sync.Mutex)

var globalBonbonLock = new(sync.Mutex)

// 實作 send message API
func handleSend(msg []byte, id int, u *user) {
	var req SendRequest
	err := json.Unmarshal(msg, &req)
	// 無法偵測出json格式是否正確
	if err != nil {
		fmt.Printf("unmarshal send cmd, %s\n", err.Error())
		return
	}

	println(req.Msg)
	u_now := time.Now()
	database.AppendMessage(id, req.Who, 0, req.Msg, u_now)
	now := u_now.UnixNano()
	ss := SendFromServer{Cmd: "sendFromServer", Who: id, Time: now, Msg: req.Msg}

	if req.Who != 0 && sendJsonToUnknownStatusID(req.Who, ss, false) == nil {
		sendJsonToOnlineID(id, respondToSend(req, now, true), false)
	} else if req.Who == 0 {
		var stranger int
		globalMatchLock.Lock()
		if stranger = u.match; stranger != -1 {
			ss.Who = 0
			sendJsonToUnknownStatusID(u.match, ss, false)
		}
		globalMatchLock.Unlock()
		if stranger == -1 {
			sendJsonToOnlineID(id, respondToSend(req, now, false), false)
		} else {
			sendJsonToOnlineID(id, respondToSend(req, now, true), false)
		}
	} else {
		sendJsonToOnlineID(id, respondToSend(req, now, false), false)
	}
}

func handleBonbon(id int, u *user) {
	fmt.Printf("%d bonbon\n", id)
	var success = false
	onlineLock.RLock()
	globalMatchLock.Lock()
	globalBonbonLock.Lock()
	var stranger *user
	strangerID := u.match
	if strangerID == -1 {
		fmt.Printf("沒有connect就bonbon\n")
		goto bonbonUnlock
	}
	stranger = onlineUser[strangerID]
	if stranger == nil {
		fmt.Printf("陌生人已經離線或不存在\n")
		goto bonbonUnlock
	}

	if stranger.bonbon == false {
		fmt.Printf("%d bonbon: 對方未bonbon\n", id)
		u.bonbon = true
	} else if stranger.bonbon == true {
		fmt.Printf("%d bonbon: 成為朋友\n", id)
		u.bonbon = false
		stranger.bonbon = false
		u.match = -1
		stranger.match = -1
		success = true
	}
bonbonUnlock:
	globalBonbonLock.Unlock()
	globalMatchLock.Unlock()
	onlineLock.RUnlock()

	sendJsonToOnlineID(id, bonbonResponse{OK: true, Cmd: "bonbon"}, false)

	if success {
		err := database.MakeFriendship(id, strangerID)
		if err != nil {
			return
		}
		strangerNick, err := database.GetNickNameOfFriendship(id, strangerID)
		if err != nil {
			return
		}
		myNick, err := database.GetNickNameOfFriendship(strangerID, id)
		if err != nil {
			return
		}
		sendJsonToOnlineID(id, NewFriendCmd{Cmd: "new_friend", Who: strangerID, Nick: strangerNick}, false)
		sendJsonToUnknownStatusID(
			strangerID,
			NewFriendCmd{Cmd: "new_friend", Who: id, Nick: myNick},
			false,
		)
	}
}

// handle account settings update
func handleUpdateSettings(msg []byte, id int) {
	// decode JSON request
	var request updateSettingsRequest
	err := json.Unmarshal(msg, &request)
	if err != nil {
		response := updateSettingsResponse{OK: false, Cmd: "setting", Setting: request.Setting}
		sendJsonToOnlineID(id, &response, false)
		return
	}

	// update database
	err = database.SetSignature(id, request.Setting.Sign)
	if err != nil {
		response := updateSettingsResponse{OK: false, Cmd: "setting", Setting: request.Setting}
		sendJsonToOnlineID(id, &response, false)
		return
	}

	// TODO 告訴所有人此人改簽名檔
	// send success response
	response := updateSettingsResponse{OK: true, Cmd: "setting", Setting: request.Setting}
	sendJsonToOnlineID(id, &response, false)
}

// XXX: 下版功能 handle setting nickname of friends
func handleSetNickName(msg []byte, id int) {
	// TODO: 修正response
	var request setNickNameRequest
	err := json.Unmarshal(msg, &request)
	if err != nil {
		response := simpleResponse{OK: false}
		sendJsonToOnlineID(id, &response, false)
		return
	}

	// update database
	err = database.SetNickNameOfFriendship(id, request.Who, request.NickName)
	if err != nil {
		response := simpleResponse{OK: false}
		sendJsonToOnlineID(id, &response, false)
		return
	}

	// send success response
	response := simpleResponse{OK: true}
	sendJsonToOnlineID(id, &response, false)
}

type requestInChannel struct {
	ID      int
	Msg     []byte
	Special string
	Conn    *websocket.Conn
	User    *user
}

var responseChannel = make(map[*websocket.Conn]chan *user)

var requestChannel = make(chan requestInChannel)

func CommandComsumer() {
	for {
		req := <-requestChannel
		id := req.ID
		conn := req.Conn
		println(req.Special)
		if req.Special == "init" {
			user, err := initOnline(id, conn)
			if err != nil {
				fmt.Printf("send initialize message to %d fail, %s\n", id, err)
			}
			responseChannel[req.Conn] <- user
			continue
		} else if req.Special == "close" {
			continue
		}

		user := req.User
		msg := req.Msg
		var decodedMsg map[string]interface{}
		json.Unmarshal(req.Msg, &decodedMsg)
		switch decodedMsg["Cmd"] {
		// NOTE: 各種cmd其實也可以僅傳入id，但傳入user可增進效能（不用再搜一次map）
		case "setting":
			handleUpdateSettings(msg, id)
			// XXX: 誤用API，此為下版功能
			// case "change_nick":
			// 	handleSetNickName(msg, id)
		case "connect":
			fmt.Printf("id %d try connect\n", id)
			handleConnect(msg, id, user)
		case "send":
			handleSend(msg, id, user)
		case "disconnect":
			handleDisconnect(id)
		case "bonbon":
			handleBonbon(id, user)
		case "history":
			handleHistory(msg, id, user)
		default:
			fmt.Println("未知的請求")
		}
	}
}

// ChatHandler 一個gin handler，為websocket之入口
func ChatHandler(id int, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("establish connection, %s\n", err.Error())
		return
	}
	responseChannel[conn] = make(chan *user)
	requestChannel <- requestInChannel{ID: id, Conn: conn, Special: "init"}
	// user, err := initOnline(id, conn)
	user := <-responseChannel[conn]
	// TODO 通知所有人此人上線
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("id %d: %s\n", id, msg)
		requestChannel <- requestInChannel{ID: id, User: user, Conn: conn, Msg: msg}
	}
	fmt.Printf("%d id leave")
	clearOffline(id, conn)
}
