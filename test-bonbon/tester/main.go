package main

import (
	"bonbon/communicate"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

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
	testInit()
}
