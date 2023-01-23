package models;

import (
	"fmt"
	"github.com/gorilla/websocket"
	// "strconv"
	"sync"
	// "github.com/Moody0101-X/Go_Api/DATABASE"
)

var ClientPool ClientHub;

type ClientHub struct {
	sync.RWMutex
	SocketClients map[int]Client
	Initialized   bool
}

func (Hub *ClientHub) GetClient(id int) (Client, bool) {
	Hub.Lock()
	val, ok := Hub.SocketClients[id]
	Hub.Unlock()
	return val, ok;
}

func (Hub *ClientHub) AddClient(Addr string, Uuid int, Conn *websocket.Conn) (*Client, bool) {
	
	if(!Hub.Initialized) {
		Hub.SocketClients = make(map[int]Client)
		Hub.Initialized = true;
	}

	var New Client
	New.Addr = Addr
	New.Uuid = Uuid
	New.Conn = Conn
	Hub.Lock();
	Hub.SocketClients[Uuid] = New
	Hub.Unlock();
	New.logClient()
	val, ok := Hub.GetClient(Uuid);
	return &val, ok
}

func (Hub *ClientHub) BroadCastJSON(Msg interface{}, skip int) {
	// Send the message for everyone..

	if(!Hub.Initialized) {
		return
	}

	Hub.Lock();

   // Sends The action to all people in
	{
		for id, Client_ := range Hub.SocketClients { 
			if(id != skip) {
				Client_.SendJSON(Msg);
			}
		}
	}

	Hub.Unlock();
}

type Client struct {
	Addr string
	Uuid int
	Conn *websocket.Conn
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