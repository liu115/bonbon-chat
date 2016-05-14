package communicate

import (
	"bonbon/database"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	// "sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type user struct {
	conns     []*websocket.Conn
	match     int // -1表示無配對
	matchType string
	bonbon    bool
}

var onlineUser = make(map[int]*user)

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
	now := u_now.UnixNano()
	ss := SendFromServer{Cmd: "sendFromServer", Who: id, Time: strconv.FormatInt(now, 10), Msg: req.Msg}

	if req.Who != 0 {
		err = database.AppendMessage(id, req.Who, 0, req.Msg, u_now)
		if err != nil {
			fmt.Printf("AppendMessage fail, %s\n", err.Error())
			sendJsonToOnlineID(id, respondToSend(req, strconv.FormatInt(now, 10), false))
		} else {
			sendJsonToUnknownStatusID(req.Who, ss)
			sendJsonToOnlineID(id, respondToSend(req, strconv.FormatInt(now, 10), true))
		}
	} else if req.Who == 0 {
		var stranger int
		if stranger = u.match; stranger != -1 {
			ss.Who = 0
			sendJsonToUnknownStatusID(u.match, ss)
		}
		if stranger == -1 {
			sendJsonToOnlineID(id, respondToSend(req, strconv.FormatInt(now, 10), false))
		} else {
			sendJsonToOnlineID(id, respondToSend(req, strconv.FormatInt(now, 10), true))
		}
	} else {
		sendJsonToOnlineID(id, respondToSend(req, strconv.FormatInt(now, 10), false))
	}
}

func handleBonbon(id int, u *user) {
	fmt.Printf("%d bonbon\n", id)
	var stranger *user
	strangerID := u.match
	fmt.Printf("strangerID is %s\n", strangerID)
	if strangerID == -1 {
		fmt.Printf("沒有connect就bonbon\n")
		sendJsonToOnlineID(id, BonbonResponse{OK: false, Cmd: "bonbon"})
		return
	}
	stranger = onlineUser[strangerID]
	if stranger == nil {
		fmt.Printf("陌生人已經離線或不存在\n")
		sendJsonToOnlineID(id, BonbonResponse{OK: false, Cmd: "bonbon"})
		return
	}

	if stranger.bonbon == false {
		fmt.Printf("%d bonbon: 對方未bonbon\n", id)
		u.bonbon = true
		sendJsonToOnlineID(id, BonbonResponse{OK: true, Cmd: "bonbon"})
		return
	} else if stranger.bonbon == true {
		fmt.Printf("%d bonbon: 成為朋友\n", id)
		u.bonbon = false
		stranger.bonbon = false
		u.match = -1
		stranger.match = -1

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
		sendJsonToOnlineID(id, BonbonResponse{OK: true, Cmd: "bonbon"})
		sendJsonToOnlineID(id, NewFriendCmd{Cmd: "new_friend", Who: strangerID, Nick: strangerNick})
		sendJsonToUnknownStatusID(
			strangerID,
			NewFriendCmd{Cmd: "new_friend", Who: id, Nick: myNick},
		)
	}
}

// handle account settings update
func handleUpdateSettings(msg []byte, id int) {
	// decode JSON request
	var request UpdateSettingsRequest
	err := json.Unmarshal(msg, &request)
	if err != nil {
		response := UpdateSettingsResponse{OK: false, Cmd: "setting", Setting: request.Setting}
		sendJsonToOnlineID(id, &response)
		return
	}

	// update database
	err = database.SetSignature(id, request.Setting.Sign)
	if err != nil {
		response := UpdateSettingsResponse{OK: false, Cmd: "setting", Setting: request.Setting}
		sendJsonToOnlineID(id, &response)
		return
	}

	// 通知朋友
	friendships, err := database.GetFriendships(id)
	if err == nil {
		for i := 0; i < len(friendships); i++ {
			fmt.Printf("%d 通知 %d 換簽名檔\n", id, friendships[i].FriendID)
			sendJsonToUnknownStatusID(
				friendships[i].FriendID,
				SignCmd{Cmd: "change_sign", Who: id, Sign: request.Setting.Sign},
			)
		}
	} else {
		fmt.Printf("getFriendships: %s", err.Error())
	}
	// send success response
	response := UpdateSettingsResponse{OK: true, Cmd: "setting", Setting: request.Setting}
	sendJsonToOnlineID(id, &response)
}

func handleSetNickName(msg []byte, id int) {
	var request SetNickNameRequest
	err := json.Unmarshal(msg, &request)
	if err != nil {
		response := SetNickNameResponse{OK: false, Cmd: "set_nick", Who: request.Who, Nick: request.Nick}
		sendJsonToOnlineID(id, &response)
		return
	}

	// update database
	err = database.SetNickNameOfFriendship(id, request.Who, request.Nick)
	if err != nil {
		fmt.Printf("SetNickNameOfFriendship in handleSetNickName: %s\n", err.Error())
		response := SetNickNameResponse{OK: false, Cmd: "set_nick", Who: request.Who, Nick: request.Nick}
		sendJsonToOnlineID(id, &response)
		return
	}

	// send success response
	response := SetNickNameResponse{OK: true, Cmd: "set_nick", Who: request.Who, Nick: request.Nick}
	sendJsonToOnlineID(id, &response)
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
		if req.Special == "init" {
			user, err := initOnline(id, conn)
			if err != nil {
				fmt.Printf("send initialize message to %d fail, %s\n", id, err)
			}
			responseChannel[req.Conn] <- user
			continue
		} else if req.Special == "close" {
			clearOffline(id, conn)
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
		case "set_nick":
			handleSetNickName(msg, id)
		case "connect":
			fmt.Printf("id %d try connect\n", id)
			handleConnect(msg, id, user)
		case "send":
			handleSend(msg, id, user)
		case "disconnect":
			handleDisconnect(id, user)
		case "bonbon":
			handleBonbon(id, user)
		case "history":
			handleHistory(msg, id, user)
		case "read":
			handleRead(msg, id, user)
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
	user := <-responseChannel[conn]
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("id %d: %s\n", id, msg)
		requestChannel <- requestInChannel{ID: id, User: user, Conn: conn, Msg: msg}
	}
	requestChannel <- requestInChannel{ID: id, Conn: conn, Special: "close"}
}
