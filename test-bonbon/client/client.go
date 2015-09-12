package client

import (
	"bonbon/communicate"
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

type Client struct {
	Conn *websocket.Conn
}

func (c *Client) Close() {
	c.Conn.Close()
}

func (c *Client) Send(id int, s string) {
	c.Conn.WriteJSON(communicate.SendRequest{
		Cmd:   "send",
		Who:   id,
		Msg:   s,
		Order: 0,
	})
}

func CreateClient(id int) *Client {
	return &Client{Conn: createConn(id)}
}

func CreateAndReceiveInit(id int) *Client {
	c := CreateClient(id)
	_, _, err := c.Conn.ReadMessage()
	if err != nil {
		fmt.Printf("in CreateAndReceiveInit, %s", err.Error())
		return nil
	}
	return c
}
