package database;

import "github.com/Moody0101-X/Go_Api/models"
import "fmt"
/*

CREATE TABLE CONVERSATIONS (
    ID INTEGER PRIMARY KEY AUTOINCREMENT, 
    Fpair INTEGER,
	Spair INTEGER,
	timestamp Date,
);

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

*/

func CreateNewDiscussion(topic_id int, other_id int) (conversation_id int) {
	conversation_id = DiscussionExists(topic_id, other_id);
	
	if conversation_id == -1 {
		stmt, _ := dataBase.Prepare("INSERT INTO CONVERSATIONS(Fpair, Spair, timestamp) VALUES(?, ?, datetime())")
		_, err := stmt.Exec(topic_id, other_id)
		if err != nil {

			fmt.Println("An error accured while appending new discussion!")
			fmt.Println("", err.Error())

		}
	}

	conversation_id = DiscussionExists(topic_id, other_id);

	return conversation_id;
}


func DiscussionExists(topic_id int, other_id int) int {
	
	var conversation_id int = -1;

	row, err := dataBase.Query("SELECT ID FROM CONVERSATIONS WHERE Fpair=? AND Spair=?", topic_id, other_id)
	
	if err != nil {
		return conversation_id;
	}

	defer row.Close()

	for row.Next() {
		row.Scan(&conversation_id);
	}

	return conversation_id;
}

func SendMessage(client models.Client, Message models.UMessage) {
	
	Message.Topic_id = client.Uuid;	
	c, ok := models.ClientPool.GetClient(Message.Other_id);
	if ok { Message.Send(&c) }

	//TODO we add COnversation to reg
	conversation_id := CreateNewDiscussion(Message.Topic_id , Message.Other_id);
	Message.ConversationId = conversation_id;
	//TODO we add the message to db.
	stmt, _ := dataBase.Prepare("INSERT INTO MESSAGES(Msg, MsgType, Coversation_id, topic_id, other_id, ts, seen) VALUES(?, ?, ?, ?, ?, datetime(), 0)")
	_, err := stmt.Exec(Message.Data.Text, Message.Data.MsgType, Message.ConversationId, Message.Topic_id , Message.Other_id)

	if err != nil {
		fmt.Println("THERE WAS AN ERROR ADDING USER MESSAGE TO DB")
		fmt.Println(err.Error())
	}

}

// TODO :-)
// func GetDisscussionMessagesById(id_ int) []models.UMessage {}
// func GetUserDiscussions(User_id int) []models.Discussion {

// 	var Discussions []models.Discussion;

// 	row, err := dataBase.Query("SELECT * FROM CONVERSATIONS WHERE Fpair=? OR Spair=?", User_id);

// 	defer row.close();

// 	var temp models.Discussion;

// 	row.Next() {
// 		row.Scan(&temp.Id_);
// 	}	
// }


// type DataFrame struct {
// 	Text    string `json:"text"`
// 	MsgType string `json:"mt"`
// }

// type UMessage struct {
// 	Id_         			int          `json:"id"`
// 	ConversationId          int          `json:"conversation_id"`
// 	Data 					DataFrame    `json:"data"`
// 	Other_id 				int 	     `json:"other_id"`
// 	Topic_id 				int 	     `json:"topic_id"`
// 	timeStamp 				time.Time    `json:"ts"`
// }