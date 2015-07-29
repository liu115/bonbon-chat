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
	runtime.GOMAXPROCS(runtime.NumCPU())

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

	// setup server
	gin.SetMode(config.Mode)

	// set routes
	app := gin.Default()

	// add routes for debug mode
	if config.Mode == "debug" {
		app.GET("/test/chat/:id", func(c *gin.Context) {
			idStr := c.Param("id")
			id, err := strconv.Atoi(idStr)
			if err == nil {
				communicate.ChatHandler(id, c)
			} else {
				c.String(404, "not found")
			}
		})

		app.GET("/test/create-account-by-token/:token", func(c *gin.Context) {
			token := c.Param(":token")
			account, err := database.CreateAccountByToken(token)
			if err != nil {
				c.String(404, err.Error())
			}

			c.String(200, strconv.Itoa(account.ID))
		})

		app.GET("/test/make-friendship/:id1/:id2", func(c *gin.Context) {
			id1, err := strconv.Atoi(c.Param("id1"))
			if err != nil {
				c.String(404, err.Error())
			}

			id2, err := strconv.Atoi(c.Param("id2"))
			if err != nil {
				c.String(404, err.Error())
			}

			if id1 < 1 || id2 < 1 {
				c.String(404, "illegal id")
			}

			err = database.MakeFriendship(id1, id2)
			if err != nil {
				c.String(404, err.Error())
			}

			c.String(200, "success")
		})

		app.GET("/test/remove-friendship/:id1/:id2", func(c *gin.Context) {
			id1, err := strconv.Atoi(c.Param("id1"))
			if err != nil {
				c.String(404, err.Error())
			}

			id2, err := strconv.Atoi(c.Param("id2"))
			if err != nil {
				c.String(404, err.Error())
			}

			if id1 < 1 || id2 < 1 {
				c.String(404, "illegal id")
			}

			err = database.RemoveFriendship(id1, id2)
			if err != nil {
				c.String(404, err.Error())
			}

			c.String(200, "success")
		})
	}

	// routes for production puropose
	app.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "./static/chat.html")
	})
	app.Static("/static/", "./static")

	// run server
	app.Run(config.Address)
}
