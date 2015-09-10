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

func describe(d string) {
	println(d)
}

func judge(b bool, describe string) {
	if b {
		color.Green("    ✓ " + describe)
	} else {
		color.Red("    ✗ " + describe)
	}
}

var signatures = [...]string{
	"",
	"我的征途是星辰大海",
	"幹0糧母豬滾喇幹",
	"當一個人真的渴望某樣東西時，他就會無恥到某個程度",
	"男人要死，就死在選總統的路上",
	"讓你難過的事情，有一天，你一定會笑著說出來",
}

var testsuite = [...]func(){
	func() {
		describe(`
一個使用者登入
測試API: init
		`)
		clearDB()

		createAccount(1, signatures[1])

		conn := createConn(1)
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
		var req communicate.InitCmd
		json.Unmarshal(msg, &req)

		ok := true
		ok = ok && req.Cmd == "init"
		ok = ok && req.Setting.Sign == signatures[1]
		judge(ok, "初始回傳正確Cmd、簽名檔")
	},
	func() {
		describe(`
兩個互為好友者登入
測試API: init, status
		`)

		clearDB()

		createAccount(1, signatures[1])
		createAccount(2, signatures[2])
		database.MakeFriendship(1, 2)

		conn := createConn(1)
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
		var req communicate.InitCmd
		json.Unmarshal(msg, &req)
		ok := true
		ok = ok && (req.Friends[0].ID == 2)
		ok = ok && (req.Friends[0].Sign == signatures[2])
		ok = ok && (req.Friends[0].Status == "off")
		judge(ok, "id1回傳正確朋友名單及狀態")

		conn2 := createConn(2)
		_, msg, err = conn2.ReadMessage()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
		json.Unmarshal(msg, &req)
		ok = true
		ok = ok && (req.Friends[0].ID == 1)
		ok = ok && (req.Friends[0].Sign == signatures[1])
		ok = ok && (req.Friends[0].Status == "on")
		judge(ok, "id2回傳正確朋友名單及狀態")

		_, msg, err = conn.ReadMessage()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
		var statusReq communicate.StatusCmd
		json.Unmarshal(msg, &statusReq)
		ok = true
		ok = ok && (statusReq == communicate.StatusCmd{
			Cmd:    "status",
			Who:    2,
			Status: "on",
		})
		judge(ok, "id2上線主動通知")
	},
}

func main() {
	for _, test := range testsuite {
		test()
	}
}
