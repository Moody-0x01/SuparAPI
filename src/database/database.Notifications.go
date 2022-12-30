package database;

import (
	"fmt"
	"io"
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
	return Nots;
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

func SetSeenForNotification(uuid int, notificationId_ int) {
	// TODO: set Seen to true (1)
	stmt, _ := dataBase.Prepare("UPDATE NOTIFICATIONS SET Seen=1 WHERE UUID=? AND ID=?")

	_, err := stmt.Exec(uuid, notificationId_)
	if err != nil {
		fmt.Println("db err: ", err)
	}
}

func pushNotificationForUser(NotificaionObject models.Notification, suffixTxt string) {
	// TODO: add the notification to db... having a prob heree..

	Client, ok := models.SocketClients[NotificaionObject.Uuid];

	fmt.Println("NOTE OBJECT: ", NotificaionObject)

    if ok {
		NotificaionObject.User_ = GetUserById(NotificaionObject.Actorid);
		NotificaionObject.Text = NotificaionObject.User_.UserName + " " + suffixTxt;
    	Client.SendJSON(NotificaionObject)
    }

    AddNewNotification(NotificaionObject);
}

func HandleConnextionForNotifications(c *models.Client) {
	for {
		var seenFlag models.NotificationSeenFlag
		err := c.Conn.ReadJSON(&seenFlag)

		if err != nil {
			if err == io.EOF {
				fmt.Printf("User with id: %s DISCONNECTED.\n", c.Uuid);
				break;
			}

			fmt.Println(err);

			continue;
		}

		seenFlag.LogSeen()
		
		// TODO: Save the specified notification new state => seen
		// TODO: Send the success message.
		SetSeenForNotification(seenFlag.Uuid, seenFlag.Id_);
	}
}
