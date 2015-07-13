package main

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

func randomHistory() []byte {
	msg := [5]string{"11111", "2222", "3333", "444444", "5555555"}
	var j = map[string]interface{}{
		"msg": msg,
		"fun": "history",
	}
	jj, err := json.Marshal(j)
	if err != nil {
		fmt.Println(err)
	}
	return jj
}

func main() {
	app := gin.Default()
	app.GET("/chat", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err == nil {
			for {
				_, p, _ := conn.ReadMessage()
				var decodeP map[string]string
				json.Unmarshal(p, &decodeP)
				fmt.Print(decodeP)
				switch decodeP["fun"] {
				case "history":
					conn.WriteMessage(websocket.TextMessage, randomHistory())
				case "send":
					conn.WriteMessage(websocket.TextMessage, p)
				}
			}
		} else {
			fmt.Print("connect fail")
		}
	})
	app.Static("/static", "static")
	app.Run(":8080")
}
