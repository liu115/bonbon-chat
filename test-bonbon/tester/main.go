package main

import (
	"bonbon/communicate"
	"bonbon/config"
	"bonbon/database"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"test-bonbon/client"
)

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
		fmt.Printf("start\n")
		clearDB()

		createAccount(1, signatures[1])
		fmt.Printf("initialize DB end\n")
		c := client.CreateClient(1)
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
		fmt.Printf("%s", msg)
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

		clients := [3]*client.Client{nil, client.CreateClient(1)}
		_, msg, err := clients[1].Conn.ReadMessage()
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

		clients[2] = client.CreateClient(2)
		_, msg, err = clients[2].Conn.ReadMessage()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
		json.Unmarshal(msg, &req)
		ok = true
		ok = ok && (req.Friends[0].ID == 1)
		ok = ok && (req.Friends[0].Sign == signatures[1])
		ok = ok && (req.Friends[0].Status == "on")
		judge(ok, "id2回傳正確朋友名單及狀態")

		_, msg, err = clients[1].Conn.ReadMessage()
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

		clients[2].Close()
		_, msg, err = clients[1].Conn.ReadMessage()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
		json.Unmarshal(msg, &statusReq)
		ok = true
		ok = ok && (statusReq == communicate.StatusCmd{
			Cmd:    "status",
			Who:    2,
			Status: "off",
		})
		judge(ok, "id2下線主動通知")
	},
	func() {
		describe(`
兩人登入，測試互傳訊息
測試API: send
			`)
		clearDB()
		createAccount(1, signatures[1])
		createAccount(2, signatures[2])
		database.MakeFriendship(1, 2)
		clients := [...]*client.Client{nil, client.CreateAndReceiveInit(1), client.CreateAndReceiveInit(2)}
		message := "QQ"

		clients[1].Send(2, message)
		for {
			_, msg, err := clients[2].Conn.ReadMessage()
			if err != nil {
				fmt.Printf("%s", err.Error())
			}
			var j communicate.SendFromServer
			json.Unmarshal(msg, &j)
			if j.Cmd == "sendFromServer" && j.Who == 1 && j.Msg == message {
				judge(true, "id2收到id1之消息")
				break
			}
		}
	},
	func() {
		describe(`
兩非朋友登入，透過connect連線並互傳
測試API: connect, send
		`)
		clearDB()
		createAccount(1, signatures[1])
		createAccount(2, signatures[2])
		clients := [...]*client.Client{nil, client.CreateAndReceiveInit(1), client.CreateAndReceiveInit(2)}
		clients[1].Connect("stranger")
		clients[2].Connect("stranger")
		clients[1].WaitForConnected()
		clients[2].WaitForConnected()
		message := "QQ"
		clients[1].SendToStranger(message)
		for {
			_, msg, err := clients[2].Conn.ReadMessage()
			if err != nil {
				fmt.Printf("%s", err.Error())
			}
			var j communicate.SendFromServer
			json.Unmarshal(msg, &j)
			if j.Cmd == "sendFromServer" && j.Who == 0 && j.Msg == message {
				judge(true, "id2收到陌生人之消息")
				break
			}
		}
	},
	func() {
		describe(`
兩非朋友登入，以bonbon成為朋友
測試API: bonbon
		`)
		clearDB()
		createAccount(1, signatures[1])
		createAccount(2, signatures[2])
		clients := [...]*client.Client{nil, client.CreateAndReceiveInit(1), client.CreateAndReceiveInit(2)}
		clients[1].Connect("stranger")
		clients[2].Connect("stranger")
		clients[1].WaitForConnected()
		clients[2].WaitForConnected()
		clients[1].Bonbon()
		clients[2].Bonbon()
		for {
			_, msg, err := clients[1].Conn.ReadMessage()
			if err != nil {
				fmt.Printf("%s", err.Error())
			}
			var j communicate.NewFriendCmd
			json.Unmarshal(msg, &j)
			if j.Cmd == "new_friend" && j.Who == 2 { // TODO: 檢測j.Nick
				judge(true, "id1收到new_friend命令")
				break
			}
		}
		for {
			_, msg, err := clients[2].Conn.ReadMessage()
			if err != nil {
				fmt.Printf("%s", err.Error())
			}
			var j communicate.NewFriendCmd
			json.Unmarshal(msg, &j)
			if j.Cmd == "new_friend" && j.Who == 1 { // TODO: 檢測j.Nick
				judge(true, "id2收到new_friend命令")
				break
			}
		}
		judge(checkFriendship(1, 2), "兩人在資料庫中已為好友")
	},
	func() {
		describe(`
兩朋友登入，互傳後檢查歷史
測試API: history
		`)
		clearDB()
		createAccount(1, signatures[1])
		createAccount(2, signatures[2])
		database.MakeFriendship(1, 2)
		clients := [...]*client.Client{nil, client.CreateAndReceiveInit(1), client.CreateAndReceiveInit(2)}
		// TODO: 以亂數的字串、數量去測試
		msgs := [...]communicate.HistoryMsg{
			communicate.HistoryMsg{Text: "hiiii", From: 1},
			communicate.HistoryMsg{Text: "xxxxx", From: 2},
			communicate.HistoryMsg{Text: "xxxxx", From: 2},
			communicate.HistoryMsg{Text: "xxxxx", From: 2},
			communicate.HistoryMsg{Text: "ifjso", From: 1},
		}
		for _, msg := range msgs {
			var to int
			if msg.From == 1 {
				to = 2
			} else {
				to = 1
			}
			clients[msg.From].Send(to, msg.Text)
			_, _, _ = clients[2].Conn.ReadMessage()
			_, _, _ = clients[1].Conn.ReadMessage()
		}
		length := len(msgs)
		clients[1].GetHistory(2, length, 3*1e18)
		_, msg, _ := clients[1].Conn.ReadMessage()
		fmt.Printf("%s", msg)
		var j communicate.HistoryResponse
		json.Unmarshal(msg, &j)
		ok := true
		for i, msg := range j.Msgs {
			th := len(msgs) - 1 - i
			if msg.Text != msgs[th].Text || msg.From != msgs[th].From {
				fmt.Printf("%s\n%s\n", msg.Text, msgs[th].Text)
				ok = false
			}
		}
		judge(ok, "成功返回正確歷史訊息")
	},
	func() {
		describe(`
壓力測試
		`)
		// clearDB()
		// createAccount(1, signatures[1])
		// createAccount(2, signatures[2])
		// database.MakeFriendship(1, 2)
		// clients := [...]*client.Client{nil, client.CreateAndReceiveInit(1), client.CreateAndReceiveInit(2)}
	},
}

func main() {
	var configPath = flag.String("config", "bonbon-develop.conf", "the path of server configuration file")
	flag.Parse()

	// load config file
	err := config.LoadConfigFile(*configPath)
	if err != nil {
		fmt.Errorf("error: cannot load config file \"%v\"", configPath)
		panic(err.Error())
	}

	for _, test := range testsuite {
		test()
	}
}
