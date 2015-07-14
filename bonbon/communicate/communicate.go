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
		for {
			_, msg, err := conn.ReadMessage()
			if err == nil {
				var decodeP map[string]string
				json.Unmarshal(msg, &decodeP)
				switch decodeP["fun"] {
				case "send":
					fmt.Println("使用者傳送訊息")
				case "history":
					fmt.Println("要求歷史")
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
