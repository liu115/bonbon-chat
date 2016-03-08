package main

import (
	"bonbon/communicate"
	"bonbon/config"
	"bonbon/database"
	"bonbon/test"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"runtime"
)

func main() {
	log.SetPrefix("[bonbon] ")
	runtime.GOMAXPROCS(runtime.NumCPU())

	// parse arguments
	var configPath = flag.String("config", "bonbon.conf", "the path of server configuration file")
	flag.Parse()

	// load config file
	err := config.LoadConfigFile(*configPath)
	if err != nil {
		log.Printf("error: cannot load config file \"%v\"", configPath)
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
		app.GET("/test/update-facebook-friends/:id", test.HandleTestUpdateFacebookFriends)
		app.GET("/test/get-facebook-friends/:id", test.HandleTestGetFacebookFriends)
		app.GET("/test/get-facebook-friends-of-friends/:id/:degree", test.HandleTestGetFacebookFriendsOfFriends)
	}

	app.GET("/chat/:token", HandleWebsocket)

	// routes for production puropose
	app.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "./static/chat.html")
	})
	app.Static("/static/", "./static")

	// run consumer
	go communicate.CommandComsumer()
	go communicate.MatchConsumer()

	// run server
	app.Run(config.Address)
}
