package main

import (
	"bonbon/account"
	"bonbon/communicate"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	app := gin.Default()
	app.GET("/init", communicate.InitHandler)
	app.GET("/chat", communicate.ChatHandler)
	app.POST("/signup", account.SignUpHandler)
	app.POST("/login", account.LoginHandler)
	app.POST("/logout", account.LogoutHandler)
	app.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "./static/chat.html")
	})
	app.Static("/static/", "./static")
	app.Run(":8080")
}
