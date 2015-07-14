package main

import (
	"bonbon/communicate"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	app := gin.Default()
	app.GET("/chat", communicate.Start)
	app.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "static/chat.html")
	})
	app.Static("/static/", "./static")
	app.Run(":8080")
}
