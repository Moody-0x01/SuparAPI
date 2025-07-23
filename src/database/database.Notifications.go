package database;

import (
	"fmt"
	// "io"
	"github.com/Moody0101-X/Go_Api/models"
)


func GetAllNotifications(ID int) []models.Notification {
	var Nots []models.Notification;
	// TODO: Get all nots from the DATABASE to be shipped to the fron-end.

	row, err := DATABASE.Query("SELECT ID, TEXT, TYPE, USER_ID, ACTOR_ID, SEEN, POST_ID, LINK FROM NOTIFICATIONS WHERE ID=? ORDER BY ID DESC", ID)
	defer row.Close()
	
	if err != nil {
		fmt.Println(err)
		return Nots
	}

	var temp models.Notification

	for row.Next() {
		row.Scan(&temp.Id_, &temp.Text, &temp.Type, &temp.Uuid, &temp.Actorid, &temp.Seen, &temp.Post_id, &temp.Link)
		temp.User_ = GetUserById(temp.Actorid);
		Nots = append(Nots, temp)
	}
	return Nots;
}

func AddNewNotification(Entry models.Notification) {
	// TODO: Add a new Notification using a notification structure.
	stmt, _ := DATABASE.Prepare("INSERT INTO NOTIFICATIONS(TEXT, TYPE, USER_ID, ACTOR_ID, SEEN, POST_ID, LINK) VALUES(?, ?, ?, ?, ?, ?, ?)")
	_, err := stmt.Exec(Entry.Text, Entry.Type, Entry.Uuid, Entry.Actorid, Entry.Seen, Entry.Post_id, Entry.Link)
	
	if err != nil {
		fmt.Println("ERROR: ");
		fmt.Println("", err);
	}
}

func SetSeenForNotification(id int) {
	// TODO: set Seen to true (1)
	stmt, _ := DATABASE.Prepare("UPDATE NOTIFICATIONS SET SEEN=1 WHERE ID=?")

	_, err := stmt.Exec(id)
	
	if err != nil {
		fmt.Println("There was an error while adding the notification seen flag")
		fmt.Println(err)
	}
}

func pushNotificationForUser(NotificaionObject models.Notification, suffixTxt string) {
	// TODO: add the notification to db... having a prob heree..

	if(!(NotificaionObject.Actorid == NotificaionObject.Uuid)) {
		Client, ok := models.ClientPool.GetClient(NotificaionObject.Uuid);
		NotificaionObject.User_ = GetUserById(NotificaionObject.Actorid);
		NotificaionObject.Text = NotificaionObject.User_.UserName + " " + suffixTxt;

		if ok {
			var resp models.SocketMessage = NotificaionObject.EncodeToSocketResponse();
			Client.SendJSON(resp);
		}

		AddNewNotification(NotificaionObject);
	}
}

// type socketResp struct {
// 	Code int `json:"code"`
// 	Data interface{} `json:"data"`
// }
