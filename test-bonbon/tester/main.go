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

func testInit() {
	conn := createConn(1)
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	var req communicate.InitCmd
	json.Unmarshal(msg, &req)
	if req.Cmd == "init" {
		color.Green("✓ 初始回傳Cmd: init")
	} else {
		color.Red("✗ 初始回傳Cmd: init")
	}
}

func main() {
	testInit()
}
