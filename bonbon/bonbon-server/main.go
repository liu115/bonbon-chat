package main

import (
	// "encoding/json"
	// "fmt"
	"bonbon/communicate"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	communicate.Show()
	app := gin.Default()
	// app.GET("/chat", func(c *gin.Context) {
	// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	// 	if err == nil {
	// 		for {
	// 			_, p, _ := conn.ReadMessage()
	// 			var decodeP map[string]string
	// 			json.Unmarshal(p, &decodeP)
	// 			fmt.Print(decodeP)
	// 			switch decodeP["fun"] {
	// 			case "send":
	// 				conn.WriteMessage(websocket.TextMessage, p)
	// 			}
	// 		}
	// 	} else {
	// 		fmt.Print("connect fail")
	// 	}
	// })
	app.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "static/chat.html")
	})
	app.Static("/static/", "./static")
	app.Run(":8080")
}
