package communicate

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// ChatHandler 一個gin handler，為websocket之入口
func ChatHandler(id int, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err == nil {
		fmt.Printf("id %d login\n", id)
		// TODO 通知所有人此人上線
		// TODO 送出初始化訊息
		for {
			_, msg, err := conn.ReadMessage()
			if err == nil {
				var decodedMsg map[string]interface{}
				json.Unmarshal(msg, &decodedMsg)
				switch decodedMsg["cmd"] {
				case "init":
					// TODO: 需要資料庫
				case "setting":
					// TODO: 需要資料庫
				case "change_nick":
					// TODO: 需要資料庫
				case "connect":
				case "connected":
				case "send":
				case "disconnect":
				case "disconnected":
				case "new_friend":
					// TODO: 需要資料庫
				default:
					fmt.Println("未知的請求")
				}
			} else {
				fmt.Println("can't read message, client close")
				break
			}
		}
	} else {
		fmt.Printf("establish connection error: %s\n", err.Error())
	}
}
