package database;

import (
	"fmt"
	"github.com/Moody0101-X/Go_Api/models"
)

func GetAllNotifications(uuid int) []models.Notification {
	var Nots []models.Notification;
	// TODO: Get all nots from the database to be shipped to the fron-end.

	row, err := dataBase.Query("SELECT ID, Text, TYPE, UUID, ACTORID, Seen, Post_id, Link FROM NOTIFICATIONS WHERE UUID=? ORDER BY ID DESC", uuid)

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

	return Notes;
}

func AddNewNotification(Entry models.Notification) {
	// TODO: Add a new Notification using a notification structure.
	stmt, _ := dataBase.Prepare("INSERT INTO NOTIFICATIONS(Text, TYPE, UUID, ACTORID, Seen, Post_id, Link) VALUES(?, ?, ?, ?, ?, ?, ?)")
	_, err := stmt.Exec(Entry.Text, Entry.Type, Entry.Uuid, Entry.Actorid, Entry.Seen, Entry.Post_id, Entry.Link)
	
	if err != nil {
		fmt.Println("ERROR: ");
		fmt.Println("", err);
	}
}

func SetSeenForNotification(uuid int)  {
	// TODO: set Seen to true (1)
	
}

func pushNotificationForUser(NotificaionObject models.Notification, suffixTxt string) {
	// TODO: Look for the connection in the online clients object. done
	// TODO: send the Notification if the user is connected. done
	// TODO: add the notification to db...
	
	client, ok := models.SocketClients[NotificaionObject.Uuid];

    if ok {
		NotificaionObject.User_ = GetUserById(NotificaionObject.Actorid);
		Notification.Text = NotificaionObject.User_.UserName + " " + suffixTxt;
    	client.SendJSON(NotificaionObject)
    } else {
    	fmt.Println("User is OFFLINE..!")
    }

    AddNewNotification(NotificaionObject);
}
