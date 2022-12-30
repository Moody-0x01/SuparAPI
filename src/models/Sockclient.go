package models;

import (
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
	// "github.com/Moody0101-X/Go_Api/database"
)

var SocketClients = make(map[int]Client)

type Client struct {
	Addr string
	Uuid int
	Conn *websocket.Conn
}

func NewClient(Addr string, Uuid int, Conn *websocket.Conn) (*Client, bool) {
	var New Client
	New.Addr = Addr
	New.Uuid = Uuid
	New.Conn = Conn
	// Register user by id.
	SocketClients[Uuid] = New
	New.logClient()
	val, ok := SocketClients[Uuid]
	return &val, ok
}

func (c *Client) logClient() {
	fmt.Println("addr: ", c.Addr)
	fmt.Println("uuid: ", c.Uuid)
}

func (c *Client) sendMessage(msg string) (err error) {
	var conn = c.Conn;
	NewMsg := []byte(msg);
	err = conn.WriteMessage(websocket.TextMessage, NewMsg)
	return err
}

func (c *Client) SendJSON(v interface{}) (err error) {
	var conn = c.Conn;
	err = conn.WriteJSON(v)
	return err
}

func BroadCast(msg []byte, c *Client)  {
	defer c.Conn.Close();

	for id, Client_ := range SocketClients {
		
		if Client_ != *c {				
			var NewMsg string = strconv.Itoa(c.Uuid) + " said: " + string(msg)
			
			err := Client_.sendMessage(NewMsg)
			
			if err != nil {
				fmt.Println("Erorr sending message to user #", id, " :", err)
			}
		}

	}
}