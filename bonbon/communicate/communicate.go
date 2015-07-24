package communicate

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ChatHandler 一個gin handler，為websocket之入口
func ChatHandler(c *gin.Context) {
	fmt.Println("got request")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err == nil {
		// TODO: 送出初始化訊息
		for {
			_, msg, err := conn.ReadMessage()
			if err == nil {
				var decodedMsg map[string]string
				json.Unmarshal(msg, &decodedMsg)
				switch decodedMsg["cmd"] {
				case "init":
				case "setting":
				case "change_nick":
				case "connect":
				case "connected":
				case "send":
					fmt.Println("使用者傳送訊息")
				case "disconnect":
				case "disconnected":
				case "status":
				case "new_friend":
				default:
					fmt.Println("未知的請求")

				}
			} else {
				fmt.Println("can't read message, client close")
				break
			}
		}
	} else {
		fmt.Print("establish connection error")
	}
}
