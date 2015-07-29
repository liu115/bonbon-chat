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
	conns []*websocket.Conn
	match int         // -1表示無配對
	lock  *sync.Mutex // lock住conns
}

var onlineUser = make(map[int]*user)
var onlineLock = new(sync.RWMutex)

// 粒度高，將降低效能
// TODO: 改為讀寫鎖或拆分粒度
var globalMatchLock = new(sync.Mutex)

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
func handleSend(msg []byte, id int) {
	var req SendCmd
	err := json.Unmarshal(msg, &req)
	// 無法偵測出json格式是否正確
	if err != nil {
		fmt.Printf("unmarshal send cmd, %s\n", err.Error())
		return
	}
	u, err := getUserByID(id)
	if err != nil {
		fmt.Printf("handleSend, %s", err.Error())
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
func handleConnect(msg []byte, id int) {
	fmt.Printf("start handle Connect\n")
	var req connectCmd
	err := json.Unmarshal(msg, &req)
	if err == nil {
		fmt.Printf("Try choose stranger\n")
		sendJsonByID(id, connectCmdResponse{OK: true, Cmd: "connect"})
		var stranger = -1
		switch req.Type {
		case "stranger":
			StrangerLock.Lock()
			fmt.Printf("Try disconnect\n")
			disconnectByID(id)
			fmt.Printf("disconnect OK\n")
			// 此時必為無配對
			if waitingStranger == -1 || waitingStranger == id {
				waitingStranger = id
			} else {
				stranger = waitingStranger
				waitingStranger = -1
				globalMatchLock.Lock()
				onlineUser[id].match = stranger
				onlineUser[stranger].match = id
				globalMatchLock.Unlock()
			}
			StrangerLock.Unlock()
			// TODO: 取得對方的簽名
			if stranger != -1 {
				fmt.Printf("%d connect to %d\n", id, stranger)
				sendJsonByID(stranger, connectSuccess{Cmd: "connected", Sign: "XXXXXXX"})
				sendJsonByID(id, connectSuccess{Cmd: "connected", Sign: "XXXXXXX"})
			}

		case "L1_FB_friend":
		case "L2_FB_friend":
		}
	} else {
		fmt.Println("unmarshal connect cmd, %s", err.Error())
	}
}

// handle account settings update
func handleUpdateSettings(msg []byte, id int) {
	// decode JSON request
	var request updateSettingsRequest
	err := json.Unmarshal(msg, &request)
	if err != nil {
		response := simpleResponse{OK: false}
		sendJsonByID(id, &response)
		return
	}

	// update database
	err = database.SetSignature(id, request.Settings.Signature)
	if err != nil {
		response := simpleResponse{OK: false}
		sendJsonByID(id, &response)
		return
	}

	// send success response
	response := simpleResponse{OK: true}
	sendJsonByID(id, &response)
}

// handle setting nickname of friends
func handleSetNickName(msg []byte, id int) {
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

func handleBonbon(id int, strangerID int) {
	// TODO
	// err := database.MakeFriendship(id, strangerID)
	// ...
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

func initOnline(id int, conn *websocket.Conn) *user {
	// TODO: 先檢測是否存在於資料庫
	fmt.Printf("id %d login\n", id)
	onlineLock.Lock()
	if onlineUser[id] == nil {
		onlineUser[id] = &user{
			match: -1,
			conns: []*websocket.Conn{conn},
			lock:  new(sync.Mutex),
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
	// 若還在等待陌生人
	removeFromStrangerQueue(id)
	// 若還在連線
	disconnectByID(id)

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
	initOnline(id, conn)
	err = sendInitMsg(id)
	if err != nil {
		fmt.Printf("send initialize message to %d fail, %s\n", id, err)
		goto clear
	}
	// TODO 通知所有人此人上線
	for {
		_, msg, err := conn.ReadMessage()
		if err == nil {
			go func() {
				fmt.Printf("id %d: %s\n", id, msg)
				var decodedMsg map[string]interface{}
				json.Unmarshal(msg, &decodedMsg)
				switch decodedMsg["Cmd"] {
				case "setting":
					handleUpdateSettings(msg, id)
				case "change_nick":
					handleSetNickName(msg, id)
				case "connect":
					fmt.Printf("id %d try connect\n", id)
					handleConnect(msg, id)
				case "send":
					handleSend(msg, id)
				case "disconnect":
					handleDisconnect(id)
				case "bonbon":
					// TODO get stranger
					// handleBonbon(id, ...stranger id...)
				default:
					fmt.Println("未知的請求")
				}
			}()
		} else {
			break
		}
	}
clear:
	clearOffline(id, conn)
}
