package main

import (
	"bonbon/communicate"
	"bonbon/config"
	"bonbon/database"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // provide sqlite3 driver
	"net"
	"net/http"
	"net/url"
	"strconv"
)

func clearDB() error {
	db, err := gorm.Open(config.DatabaseDriver, config.DatabaseArgs)
	if err != nil {
		return fmt.Errorf("cannot connect to database %v://%v", config.DatabaseDriver, config.DatabaseArgs)
	}

	db.DropTable(&database.Account{})
	db.DropTable(&database.Friendship{})
	database.InitDatabase()
	return nil
}

func createAccount(ID int, signature string) error {
	user := database.Account{ID: ID, Signature: signature}

	db, err := database.GetDB()
	if err != nil {
		return err
	}

	db.Create(&user)
	return nil
}

func createConn(id int) *websocket.Conn {
	u, err := url.Parse("http://localhost:8080/test/chat/" + strconv.Itoa(id))
	rawConn, err := net.Dial("tcp", u.Host)
	conn, _, err := websocket.NewClient(rawConn, u, http.Header{}, 1024, 1024)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	return conn
}

func judge(b bool, describe string) {
	if b {
		color.Green("✓ " + describe)
	} else {
		color.Red("✗ " + describe)
	}
}

func testInit() {
	conn := createConn(1)
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	var req communicate.InitCmd
	json.Unmarshal(msg, &req)
	judge(req.Cmd == "init", "初始回傳Cmd: init")
}

func main() {
	clearDB()
	createAccount(1, "我們的征途是星辰大海")
	// testInit()
}
