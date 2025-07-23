/*
    CREATE TABLE NOTIFICATIONS (
        ID INTEGER PRIMARY KEY AUTOINCREMENT,
        Text TEXT,
        TYPE INTEGER,
        UUID INTEGER,
        ACTORID INTEGER,
        Seen INTEGER,
        Post_id INTEGER,
        Link TEXT
    )

*/
package models;
import "fmt"

type Notification struct {
	Id_       int      `json:"id"`
	Text      string   `json:"text"`
	Type      int      `json:"type"` 
	Uuid      int      `json:"user_id"`
	Actorid   int      `json:"actorid"`
	Seen      int      `json:"seen"`
	Post_id   int      `json:"post_id"`
	Link      string   `json:"link"`
	User_     AUser    `json:"User"`
}

type NotificationSeenFlag struct {
	Id_  int `json:"id"`
	Uuid int `json:"uuid"`
}

func (n *NotificationSeenFlag) LogSeen() {
	fmt.Println("id: ", n.Id_)
	fmt.Println("uuid: ", n.Uuid)
}

func (N *Notification) EncodeToSocketResponse() SocketMessage { 
	return MakeSocketResp(NOTIFICATION, 200, N) 
}

func NewNot(t int, uuid int, actorid int) Notification {
	var new Notification;	
	new.Type = t;
	new.Uuid = uuid;
	new.Actorid = actorid;
	new.Seen = 0;
	return new;
}



