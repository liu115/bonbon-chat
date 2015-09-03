package main

import (
	"fmt"
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

func main() {
	fmt.Println("test")
	conn := createConn(1)
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	fmt.Printf("%s\n", msg)
}
