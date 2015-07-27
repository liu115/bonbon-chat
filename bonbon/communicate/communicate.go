package communicate

import (
	"encoding/json"
	"errors"
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

// 粒度高，將降低效能
// TODO: 改為讀寫鎖或拆分粒度
var globalMatchLock = new(sync.Mutex)

// -1代表目前無人
var waitingStranger = -1
var StrangerLock = new(sync.Mutex)

// 通用性安全傳送訊息
func sendJSONTo(id int, json interface{}) error {
	onlineUser[id].lock.Lock()
	if onlineUser[id] != nil && len(onlineUser[id].conns) > 0 {
		for i := 0; i < len(onlineUser[id].conns); i++ {
			onlineUser[id].conns[i].WriteJSON(json)
		}
		onlineUser[id].lock.Unlock()
		// TODO: 判斷是否成功送出
		return nil
	}
	onlineUser[id].lock.Unlock()
	return errors.New("sendJSONTo: can't send")
}

// 實作init訊息
func sendInitMsg(id int) error {
	msg, err := getInitInfo(id)
	if err != nil {
		return err
	}
	err = sendJSONTo(id, msg)
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
	if err == nil {
		println(req.Msg)
		now := time.Now().UnixNano()
		ss := SendFromServer{Cmd: "sendFromServer", Who: id, Time: now, Msg: req.Msg}
		if req.Who != 0 && sendJSONTo(req.Who, ss) == nil {
			sendJSONTo(id, respondToSend(req, now, true))
		} else if req.Who == 0 {
			var stranger int
			globalMatchLock.Lock()
			if stranger = onlineUser[id].match; stranger != -1 {
				ss.Who = 0
				sendJSONTo(onlineUser[id].match, ss)
			}
			globalMatchLock.Unlock()
			if stranger == -1 {
				sendJSONTo(id, respondToSend(req, now, false))
			} else {
				sendJSONTo(id, respondToSend(req, now, true))
			}
		} else {
			sendJSONTo(id, respondToSend(req, now, false))
		}
	} else {
		fmt.Println("unmarshal send cmd, %s", err.Error())
	}
}

// 實作隨機連結(connect) API
func handleConnect(msg []byte, id int) {
	fmt.Printf("start handle Connect\n")
	var req connectCmd
	err := json.Unmarshal(msg, &req)
	if err == nil {
		fmt.Printf("Try choose stranger\n")
		sendJSONTo(id, connectCmdResponse{OK: true, Cmd: "connect"})
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
				sendJSONTo(stranger, connectSuccess{Cmd: "connected", Sign: "XXXXXXX"})
				sendJSONTo(id, connectSuccess{Cmd: "connected", Sign: "XXXXXXX"})
			}

		case "L1_FB_friend":
		case "L2_FB_friend":
		}
	} else {
		fmt.Println("unmarshal connect cmd, %s", err.Error())
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
	fmt.Printf("stranger is %d\n", stranger)
	if stranger > 0 {
		sendJSONTo(stranger, map[string]interface{}{"Cmd": "disconnected"})
	}
}

// 實作斷線
func handleDisconnect(id int) {
	disconnectByID(id)
	sendJSONTo(id, map[string]interface{}{"OK": true, "Cmd": "disconnect"})
}

// XXX：不釋放記憶體則會被越吃越多，但百萬人也才數十MB
// 但釋放後待解決lock消失的問題
func initOnline(id int, conn *websocket.Conn) {
	// XXX: 初始化時可能被多次重入
	// TODO: 先檢測是否存在於資料庫
	fmt.Printf("id %d login\n", id)
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
	onlineUser[id].lock.Lock()
	conns := onlineUser[id].conns
	var which int
	for i := 0; i < len(conns); i++ {
		if conn == conns[i] {
			which = i
			break
		}
	}
	onlineUser[id].conns = append(conns[:which], conns[which+1:]...)

	onlineUser[id].lock.Unlock()
	fmt.Printf("id %d 下線了\n", id)
}

// ChatHandler 一個gin handler，為websocket之入口
func ChatHandler(id int, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err == nil {
	} else {
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
					// TODO: 需要資料庫
				case "change_nick":
					// TODO: 需要資料庫
				case "connect":
					fmt.Printf("id %d try connect\n", id)
					handleConnect(msg, id)
				case "send":
					handleSend(msg, id)
				case "disconnect":
					handleDisconnect(id)
				case "bonbon":
					// TODO: 需要資料庫
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
