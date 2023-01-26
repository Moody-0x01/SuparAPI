package database;

import "github.com/Moody0101-X/Go_Api/models"
import "fmt"
/*

CREATE TABLE CONVERSATIONS (
    ID INTEGER PRIMARY KEY AUTOINCREMENT, 
    Fpair INTEGER,
	Spair INTEGER,
	timestamp Date
);

CREATE TABLE MESSAGES (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
	Msg TEXT,
	MsgType TEXT,
	Coversation_id INTEGER,
	topic_id INTEGER,
	other_id INTEGER
	ts DATE,
	seen INTEGER
);

*/

func CreateNewDiscussion(topic_id int, other_id int) (conversation_id int) {
	conversation_id = DiscussionExists(topic_id, other_id);
	
	if conversation_id == -1 {
		stmt, _ := DATABASE.Prepare("INSERT INTO CONVERSATIONS(Fpair, Spair, timestamp) VALUES(?, ?, datetime())")
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

	row, err := DATABASE.Query("SELECT ID FROM CONVERSATIONS WHERE Fpair=? AND Spair=? OR Fpair=? AND Spair=?", topic_id, other_id, other_id, topic_id)
				
	if err != nil {
		return conversation_id;
	}

	defer row.Close()

	for row.Next() {
		row.Scan(&conversation_id);
	}

	return conversation_id;
}

func SendMessage(client *models.Client, Message models.UMessage) {	
	Message.Topic_id = client.Uuid;	
	conversation_id := CreateNewDiscussion(Message.Topic_id , Message.Other_id);
	Message.ConversationId = conversation_id;
	
	c, ok := models.ClientPool.GetClient(Message.Other_id);
	
	if ok { 
		Message.Send(&c) 
	}

	//TODO we add COnversation to reg

	//TODO we add the message to db.
	// Message.Log();

	stmt, _ := DATABASE.Prepare("INSERT INTO MESSAGES(Msg, MsgType, Coversation_id, topic_id, other_id, ts, seen) VALUES(?, ?, ?, ?, ?, datetime(), 0)")
	_, err := stmt.Exec(Message.Data.Text, Message.Data.MsgType, Message.ConversationId, Message.Topic_id, Message.Other_id)
	
	if err != nil {
		fmt.Println("THERE WAS AN ERROR ADDING USER MESSAGE TO DB")
		fmt.Println(err.Error())
	}
}

func GetUserDiscussions(User_id int, Token string) models.Response {

	Discs := make(map[int]models.Discussion);
	
	/*		
		
		[
			{...},
			{...},
		] -> {
			i: {...},
			j: {...},
		}

	*/
	
	id, ok := GetUserIdByToken(Token)
	
	if ok {
		if id == User_id {
			row, err := DATABASE.Query("SELECT * FROM CONVERSATIONS WHERE Fpair=? OR Spair=? ORDER BY ID DESC", User_id, User_id);
			
			if err != nil {
				return models.MakeServerResponse(500, "Internal serevr error");
			}
			
			defer row.Close();
			
			var tempDisc models.Discussion;

			for row.Next() {
				
				// Fpair, Spair, timestamp
				row.Scan(&tempDisc.Id_, &tempDisc.Fpair, &tempDisc.Spair, &tempDisc.TimeStamp);
				tempDisc.Messages = GetMessagesByConvId(tempDisc.Id_);
				tempDisc.MessageCount = len(tempDisc.Messages);
				// A slight modification.

				Discs[tempDisc.Id_] = tempDisc;
			}

			return models.MakeGenericServerResponse(200, Discs);
		}
	}

	return models.MakeGenericServerResponse(401, "Not authorized!")
}

func GetMessagesByConvId(id int) []models.UMessage {
	
	var Messages []models.UMessage;
	row, err := DATABASE.Query("SELECT ID, Msg, MsgType, topic_id, other_id, ts FROM MESSAGES WHERE Coversation_id=? ORDER BY ID ASC", id);
	
	if err != nil {	
		fmt.Println("err in retrv user messages: ")
		fmt.Println("", err.Error())
		return Messages;
	}

	defer row.Close()

	var temp models.UMessage;
	
	for row.Next() {
		// ID, Msg, MsgType, Coversation_id, topic_id, other_id, ts
		row.Scan(&temp.Id_, &temp.Data.Text, &temp.Data.MsgType, &temp.Topic_id, &temp.Other_id, &temp.TimeStamp);
		temp.ConversationId = id;
		Messages = append(Messages, temp);
	}

	return Messages
}

func GetDiscussionById(uuid int, Token string, conversation_id int) models.Response {
	
	var Discussion models.Discussion;
	id, ok := GetUserIdByToken(Token)
	
	if ok {
		if id == uuid {
			row, err := DATABASE.Query("SELECT * FROM CONVERSATIONS WHERE ID=? ORDER BY ID ASC", conversation_id);
			if err != nil {
				fmt.Println("DB ERROR:", err);
				return models.MakeServerResponse(500, "Internal serevr error")
			}

			defer row.Close();

			for row.Next() {
				// Fpair, Spair, timestamp
				row.Scan(&Discussion.Id_, &Discussion.Fpair, &Discussion.Spair, &Discussion.TimeStamp);
				Discussion.Messages = GetMessagesByConvId(Discussion.Id_);
			}

			return models.MakeServerResponse(200, Discussion);
		}
	}

	return models.MakeServerResponse(401, "Not authorized!")
}


func MarkMessageAsSeen(id int) {

	// TODO: set Seen to true (1)
	stmt, _ := DATABASE.Prepare("UPDATE MESSAGES SET Seen=1 WHERE ID=?")
	_, err := stmt.Exec(id)

	if err != nil {
		fmt.Println("There was an error while adding the notification seen flag")
		fmt.Println(err)
	}
}