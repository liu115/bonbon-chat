package main

import (
	"strconv"
	"bonbon/communicate"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app := gin.Default()
	app.GET("/test/chat/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err == nil {
			communicate.ChatHandler(id, c)
		} else {
			c.String(404, "not found")
		}
	})
	app.GET("/init", communicate.InitHandler)
	app.POST("/login", LoginHandler)
	app.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "./static/chat.html")
	})
	app.Static("/static/", "./static")
	app.Run(":8080")
}
