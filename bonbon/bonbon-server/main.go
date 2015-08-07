package main

import (
	"bonbon/config"
	"bonbon/database"
	"bonbon/test"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
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
		app.GET("/test/chat/:id", test.HandleTestWebsocket)
		app.GET("/test/create-account-by-token/:token", test.HandleTestCreateAccountByToken)
		app.GET("/test/make-friendship/:id1/:id2", test.HandleTestMakeFriendship)
		app.GET("/test/remove-friendship/:id1/:id2", test.HandleTestRemoveFriendship)
	}

	// routes for production puropose
	app.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "./static/chat.html")
	})
	app.Static("/static/", "./static")

	// run server
	app.Run(config.Address)
}
