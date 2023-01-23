// READ CODE NOT DOCUMENTATION !

package models

import (
	"time"
	"fmt"
)

const (
	Audio = "Audio"
	Image = "Image"
	Video = "Video"
	Text  = "plain-text"
)


type ValidationStructure struct {
	Token      	string      `json:"token"`
	Uuid       	int      	`json:"uuid"`
	ConvId 		int 		`json:"conversation_id"` // IN case we want a particular convo !
}

type Discussion struct {
	Id_         int         `json:"id"`
	Fpair 		int 		`json:"fpair"`
	Spair		int 		`json:"spair"`
	Messages 	[]UMessage 	`json:"messages"`
	TimeStamp 	time.Time   `json:"ts"`
}

type DataFrame struct {
	Text    string `json:"text"`
	MsgType string `json:"mt"`
}

type UMessage struct {
	Id_         			int          `json:"id"`
	ConversationId          int          `json:"conversation_id"`
	Data 					DataFrame    `json:"data"`
	Other_id 				int 	     `json:"other_id"`
	Topic_id 				int 	     `json:"topic_id"`
	TimeStamp 				time.Time    `json:"ts"`
}

func (m *UMessage) Log() {
	fmt.Println("FROM: ", m.Topic_id);
	fmt.Println("TO: ", m.Other_id);
	fmt.Println("AT: ", time.Now());
	fmt.Println("MSG-CONTENT: ", m.Data.Text);
	fmt.Println("Length: ", len(m.Data.Text));
}

func NewMessage(df DataFrame, t_id int) *UMessage {
	return &UMessage{
		Data: df,
		Other_id: t_id,
		TimeStamp: time.Now(),
	}
}

func NewDataFrame(T string, MT string) *DataFrame {
	return	&DataFrame{
		Text: T, 
		MsgType: MT,
	}
}

func NewDiscussion(fpair int, spair int, MList []UMessage) *Discussion {
	return &Discussion{
		Fpair: fpair,
		Spair: spair,
		Messages: MList,
		TimeStamp: time.Now(),
	}
}



func (Message *UMessage) Send(c *Client) {

	SocketMsg := Message.EncodeToSocketResponse()

	err := c.SendJSON(SocketMsg)
	if err != nil {
		fmt.Println("Could not send to User: ")	
		fmt.Println(err.Error())
	}
}

func (Message *UMessage) SendToU(uuid int) {
	c, ok := ClientPool.GetClient(uuid);

	if ok { Message.Send(&c) }
}

func (Message *UMessage) EncodeToSocketResponse() SocketMessage {
	return MakeSocketResp(MSG, 200, Message);
}

/* 
plan
TODO:
	Make new tables:
		
		Messages:

		CREATE TABLE MESSAGES (
		    ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Msg TEXT,
			MsgType TEXT,
			Coversation_id fk,
			topic_id INTEGER,
			other_id INTEGER
			ts DATE,
			seen INTEGER 0,
		)
			

		CONVERSATION:
		
			CREATE TABLE CONVERSATIONS (
		        ID INTEGER PRIMARY KEY AUTOINCREMENT, 
		        Fpair INTEGER,
				Spair INTEGER,
				TimeStamp Date,
	    	);
			
			
		Messages # A json to capture all the messages that belong to this conversation!!
			

	first scenario:

	A user tries to send a msg, picks use #0000 then sends to the server this chunck
	
	msg: Hi
	msgT: text
	other_id: 0000
	ts: data...

	firstly we check if the conversation already exists, if not:
		INSERT INTO CONVERSATIONS(Fpair, Spair, ts) VALUES(topic_id, other_id, ts)
	then we send a msg via connexion using topic_id, if he is online he gets the messgae, if not:
		INSERT INTO Messages(Msg, MsgType, Conversation_id, topic_id, ts) VALUES(Msg, MsgType, Conversation_id, topic_id, ts)

	to find your conversation:
		firts query your conversation table.
		SELECT ID WHERE Fpair=MY_ID OR Spair=MY_ID
		---> ID we use it to retrieve all the messages.
		SELECT * FROM MESSAGES WHERE CONVERSATION_ID=ID 
*/