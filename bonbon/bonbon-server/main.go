package main

import (
	"bonbon/communicate"
	"bonbon/config"
	"bonbon/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"strconv"
)

func main() {
	// load config file
	err := config.LoadConfigFile("bonbon.conf")
	if err != nil {
		panic(err.Error())
	}

	// init database
	err = database.InitDatabase()
	if err != nil {
		panic(err.Error())
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	// setup server
	gin.SetMode(config.Mode)

	// set routes
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
	// app.POST("/login", LoginHandler)
	app.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "./static/chat.html")
	})
	app.Static("/static/", "./static")

	// run server
	app.Run(config.Address)
}
