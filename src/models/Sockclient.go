package models;

import (
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
)

var SocketClients = make(map[int]client)

type client struct {
	Addr string
	Uuid int
	Conn *websocket.Conn
	IsOpen bool
}


func NewClient(Addr string, Uuid int, Conn *websocket.Conn, IsOpen bool) {
	var New client
	New.Addr = Addr
	New.Uuid = Uuid
	New.Conn = Conn
	New.IsOpen = IsOpen
	// Register user by id.
	SocketClients[Uuid] = New
	New.logClient()
	go New.handleConn()
}

func (c *client) logClient() {
	fmt.Println("addr: ", c.Addr)
	fmt.Println("uuid: ", c.Uuid)
}

func (c *client) sendMessage(msg string) (err error) {
	var conn = c.Conn;
	NewMsg := []byte(msg);
	err = conn.WriteMessage(websocket.TextMessage, NewMsg)
	return err
}

func (c *client) SendJSON(v interface{}) (err error) {
	var conn = c.Conn;
	err = conn.WriteJSON(v)
	return err
}

func (c *client) handleConn() {
	for {
		_, message, err := c.Conn.ReadMessage()
		
		if err != nil {
			fmt.Println(err)
			break
		}

		BroadCast(message, c); // Broadcasting a message to every person that is subsribed to our notification socket pool.
	}
}


func BroadCast(msg []byte, c *client)  {
	defer c.Conn.Close();

	for id, client_ := range SocketClients {
		if client_ != *c {				
			var NewMsg string = strconv.Itoa(c.Uuid) + " said: " + string(msg)
			
			err := client_.sendMessage(NewMsg)
			
			if err != nil {
				fmt.Println("Erorr sending message to user #", id, " :", err)
			}
		}
	}
}