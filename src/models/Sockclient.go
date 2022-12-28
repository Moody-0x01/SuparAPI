package models;

var socketClients = make(map[int]client)

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
	socketClients[Uuid] = New
	fmt.Println("new client with id: ", New.Uuid);
	fmt.Println("address: ", New.Addr);
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

func (c *client) sendJSON(v interface{}) (err error) {
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

		defer c.Conn.Close()

		BroadCast(message, c); // Broadcasting a message to every person that is subsribed to our notification socket pool.
	}
}


func BroadCast(msg []byte, c client)  {
	for id, client_ := range socketClients {
		if client_ != *c {				
			var NewMsg string = strconv.Itoa(c.Uuid) + " said: " + string(msg)
			
			err := client_.sendMessage(NewMsg)
			
			if err != nil {
				fmt.Println("Erorr sending message to user #", id, " :", err)
			}
		}
	}
}