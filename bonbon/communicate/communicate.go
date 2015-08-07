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
	conns  []*websocket.Conn
	lock   *sync.Mutex // lock住conns
	match  int         // -1表示無配對
	bonbon bool
}

var onlineUser = make(map[int]*user)
var onlineLock = new(sync.RWMutex)

// 粒度高，將降低效能
// TODO: 改為讀寫鎖或拆分粒度
var globalMatchLock = new(sync.Mutex)

var globalBonbonLock = new(sync.Mutex)

// -1代表目前無人
var waitingStranger = -1
var StrangerLock = new(sync.Mutex)

// 實作init訊息
func sendInitMsg(id int) error {
	msg, err := getInitInfo(id)
	if err != nil {
		return err
	}
	err = sendJsonByID(id, msg)
	if err != nil {
		return err
	}
	return nil
}

// 實作 send message API
func handleSend(msg []byte, id int, u *user) {
	var req SendCmd
	err := json.Unmarshal(msg, &req)
	// 無法偵測出json格式是否正確
	if err != nil {
		fmt.Printf("unmarshal send cmd, %s\n", err.Error())
		return
	}

	println(req.Msg)
	now := time.Now().UnixNano()
	ss := SendFromServer{Cmd: "sendFromServer", Who: id, Time: now, Msg: req.Msg}

	if req.Who != 0 && sendJsonByID(req.Who, ss) == nil {
		sendJsonByID(id, respondToSend(req, now, true))
	} else if req.Who == 0 {
		var stranger int
		globalMatchLock.Lock()
		if stranger = u.match; stranger != -1 {
			ss.Who = 0
			sendJsonByID(u.match, ss)
		}
		globalMatchLock.Unlock()
		if stranger == -1 {
			sendJsonByID(id, respondToSend(req, now, false))
		} else {
			sendJsonByID(id, respondToSend(req, now, true))
		}
	} else {
		sendJsonByID(id, respondToSend(req, now, false))
	}
}

// 實作隨機連結(connect) API
func handleConnect(msg []byte, id int, u *user) {
	fmt.Printf("start handle Connect\n")
	var req connectCmd
	err := json.Unmarshal(msg, &req)
	if err != nil {
		fmt.Printf("unmarshal connect cmd, %s\n", err.Error())
		return
	}
	sendJsonByID(id, connectCmdResponse{OK: true, Cmd: "connect"})
	var stranger = -1
	switch req.Type {
	case "stranger":
		StrangerLock.Lock()
		disconnectByID(id)
		// 此時必為無配對，因為只能在strangerLock內建立match，而剛剛消除了match
		if waitingStranger == -1 || waitingStranger == id {
			waitingStranger = id
		} else {
			stranger = waitingStranger // 若在stranger為waitingStranger，則在strangerLock內不可能消失
			waitingStranger = -1
			globalMatchLock.Lock()
			u.match = stranger
			onlineUser[stranger].match = id
			globalMatchLock.Unlock()
		}
		StrangerLock.Unlock()

		// get signatures of both
		selfSignature, err := database.GetSignature(id)
		if err != nil {
			// TODO handle error
		}

		strangerSignature, err := database.GetSignature(stranger)
		if err != nil {
			// TODO handle error
		}

		if stranger != -1 {
			fmt.Printf("%d connect to %d\n", id, stranger)
			sendJsonByID(stranger, connectSuccess{Cmd: "connected", Sign: *selfSignature})
			sendJsonByUser(u, connectSuccess{Cmd: "connected", Sign: *strangerSignature})
		}

	case "L1_FB_friend":
	case "L2_FB_friend":
	}
}

func disconnectByID(id int) {
	var stranger int
	globalMatchLock.Lock()
	if stranger = onlineUser[id].match; stranger != -1 {
		onlineUser[id].match = -1
		onlineUser[stranger].match = -1
	}
	globalMatchLock.Unlock()
	// 將io取出鎖外操作
	fmt.Printf("%d disconnect with %d\n", id, stranger)
	if stranger > 0 {
		sendJsonByID(stranger, map[string]interface{}{"Cmd": "disconnected"})
	}
}

// 實作斷線
func handleDisconnect(id int) {
	disconnectByID(id)
	sendJsonByID(id, map[string]interface{}{"OK": true, "Cmd": "disconnect"})
}

func handleBonbon(id int) {
	fmt.Printf("%d bonbon\n", id)
	var success = false
	onlineLock.RLock()
	globalMatchLock.Lock()
	globalBonbonLock.Lock()
	var stranger *user
	strangerID := onlineUser[id].match
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
		onlineUser[id].bonbon = true
	} else if stranger.bonbon == true {
		fmt.Printf("%d bonbon: 成為朋友\n", id)
		success = true
	}
bonbonUnlock:
	globalBonbonLock.Unlock()
	globalMatchLock.Unlock()
	onlineLock.RUnlock()

	sendJsonByID(id, bonbonResponse{OK: true, Cmd: "bonbon"})

	if success {
		err := database.MakeFriendship(id, strangerID)
		if err != nil {
			return
		}
		strangerNick, err := database.GetSignature(strangerID)
		if err != nil {
			return
		}
		myNick, err := database.GetSignature(strangerID)
		if err != nil {
			return
		}
		sendJsonByID(id, newFriendFromServer{Cmd: "new_friend", Who: strangerID, Nick: *strangerNick})
		sendJsonByID(strangerID, newFriendFromServer{Cmd: "new_friend", Who: id, Nick: *myNick})
	}
}

// handle account settings update
func handleUpdateSettings(msg []byte, id int) {
	// decode JSON request
	var request updateSettingsRequest
	err := json.Unmarshal(msg, &request)
	if err != nil {
		response := updateSettingsResponse{OK: false, Cmd: "setting", Setting: request.Setting}
		sendJsonByID(id, &response)
		return
	}

	// update database
	err = database.SetSignature(id, request.Setting.Sign)
	if err != nil {
		response := updateSettingsResponse{OK: false, Cmd: "setting", Setting: request.Setting}
		sendJsonByID(id, &response)
		return
	}

	// TODO 告訴所有人此人改簽名檔
	// send success response
	response := updateSettingsResponse{OK: true, Cmd: "setting", Setting: request.Setting}
	sendJsonByID(id, &response)
}

// handle setting nickname of friends
func handleSetNickName(msg []byte, id int) {
	// TODO: 修正response
	// decode JSON request
	var request setNickNameRequest
	err := json.Unmarshal(msg, &request)
	if err != nil {
		response := simpleResponse{OK: false}
		sendJsonByID(id, &response)
		return
	}

	// update database
	err = database.SetNickNameOfFriendship(id, request.Who, request.NickName)
	if err != nil {
		response := simpleResponse{OK: false}
		sendJsonByID(id, &response)
		return
	}

	// send success response
	response := simpleResponse{OK: true}
	sendJsonByID(id, &response)
}

func initOnline(id int, conn *websocket.Conn) *user {
	// check if this account exists
	_, err := database.GetAccountByID(id)
	if err != nil {
		// TODO do error handling if id is illegal
	}

	fmt.Printf("id %d login\n", id)
	onlineLock.Lock()
	if onlineUser[id] == nil {
		onlineUser[id] = &user{
			match:  -1,
			conns:  []*websocket.Conn{conn},
			lock:   new(sync.Mutex),
			bonbon: false,
		}
	} else {
		onlineUser[id].lock.Lock()
		onlineUser[id].conns = append(onlineUser[id].conns, conn)
		onlineUser[id].lock.Unlock()
	}
	onlineLock.Unlock()
	return onlineUser[id] //此時必定還存在
}

func removeFromStrangerQueue(id int) {
	StrangerLock.Lock()
	if id == waitingStranger {
		waitingStranger = -1
	}
	StrangerLock.Unlock()
}

func clearOffline(id int, conn *websocket.Conn) {
	// 刪除本連線
	onlineLock.Lock()
	u := onlineUser[id]
	u.lock.Lock()
	conns := u.conns
	var which int
	for i := 0; i < len(conns); i++ {
		if conn == conns[i] {
			which = i
			break
		}
	}
	u.conns = append(conns[:which], conns[which+1:]...)
	if len(conns) == 0 {
		// 若還在等待陌生人
		removeFromStrangerQueue(id)
		// 若還在連線
		disconnectByID(id)
		delete(onlineUser, id)
	}
	u.lock.Unlock()
	onlineLock.Unlock()
	fmt.Printf("id %d 下線了\n", id)
}

// ChatHandler 一個gin handler，為websocket之入口
func ChatHandler(id int, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("establish connection, %s\n", err.Error())
		return
	}
	user := initOnline(id, conn)
	err = sendInitMsg(id)
	var wg sync.WaitGroup
	if err != nil {
		fmt.Printf("send initialize message to %d fail, %s\n", id, err)
		goto clear
	}
	// TODO 通知所有人此人上線
	for {
		_, msg, err := conn.ReadMessage()
		if err == nil {
			wg.Add(1)
			go func() {
				fmt.Printf("id %d: %s\n", id, msg)
				var decodedMsg map[string]interface{}
				json.Unmarshal(msg, &decodedMsg)
				switch decodedMsg["Cmd"] {
				// NOTE: 各種cmd其實也可以僅傳入id，但傳入user可增進效能（不用再搜一次map）
				case "setting":
					handleUpdateSettings(msg, id)
				case "change_nick":
					handleSetNickName(msg, id)
				case "connect":
					fmt.Printf("id %d try connect\n", id)
					handleConnect(msg, id, user)
				case "send":
					handleSend(msg, id, user)
				case "disconnect":
					handleDisconnect(id)
				case "bonbon":
					handleBonbon(id)
				default:
					fmt.Println("未知的請求")
				}
				wg.Done()
			}()
		} else {
			break
		}
	}
	wg.Wait()
clear:
	clearOffline(id, conn)
}
