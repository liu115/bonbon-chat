package communicate

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

// TODO:  使每個id可擁有多個websocket（開多分頁），改為Conn陣列
var onlineUser = make(map[int]*websocket.Conn)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type SendCmd struct {
	Cmd   string
	Who   int
	Msg   string
	Order int
}

type SendCmdResponse struct {
	OK    bool
	Who   int
	Cmd   string
	Time  int64
	Order int
}

type SendFromServer struct {
	Cmd  string
	Who  int
	Time int64
}

func respondToSend(req SendCmd, sender int, exist bool) (SendCmdResponse, SendFromServer) {
	now := time.Now().UnixNano()
	res := SendCmdResponse{
		OK:    exist,
		Who:   req.Who,
		Cmd:   req.Cmd,
		Time:  now,
		Order: req.Order,
	}
	ss := SendFromServer{
		Cmd:  "sendfromServer",
		Who:  sender,
		Time: now,
	}
	return res, ss
}

func handleSend(msg []byte, id int, conn *websocket.Conn) {
	var req SendCmd
	err := json.Unmarshal(msg, &req)
	// 無法偵測出json格式是否正確
	if err == nil {
		println(req.Msg)
		connToSend, ex := onlineUser[req.Who]
		res, ss := respondToSend(req, id, ex)
		if ex == false {
			conn.WriteJSON(res)
		} else {
			conn.WriteJSON(res)
			connToSend.WriteJSON(ss)
		}
	} else {
		fmt.Println("unmarshal send cmd, %s", err.Error())
	}
}

// ChatHandler 一個gin handler，為websocket之入口
func ChatHandler(id int, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err == nil {
		fmt.Printf("id %d login\n", id)
		onlineUser[id] = conn
		// TODO 通知所有人此人上線
		// TODO 送出初始化訊息
		for {
			_, msg, err := conn.ReadMessage()
			if err == nil {
				go func() {
					fmt.Printf("id %d: %s\n", id, msg)
					var decodedMsg map[string]interface{}
					json.Unmarshal(msg, &decodedMsg)
					switch decodedMsg["Cmd"] {
					case "init":
						// TODO: 需要資料庫
					case "setting":
						// TODO: 需要資料庫
					case "change_nick":
						// TODO: 需要資料庫
					case "connect":
					case "connected":
					case "send":
						handleSend(msg, id, conn)
					case "disconnect":
					case "disconnected":
					case "new_friend":
						// TODO: 需要資料庫
					default:
						fmt.Println("未知的請求")
					}
				}()
			} else {
				fmt.Println("can't read message, client close")
				break
			}
		}
	} else {
		fmt.Printf("establish connection, %s\n", err.Error())
	}
}
