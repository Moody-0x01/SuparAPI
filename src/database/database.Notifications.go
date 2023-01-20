package database;

import (
	"fmt"
	// "io"
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

func SetSeenForNotification(id int) {
	// TODO: set Seen to true (1)
	stmt, _ := dataBase.Prepare("UPDATE NOTIFICATIONS SET Seen=1 WHERE ID=?")

	_, err := stmt.Exec(id)
	
	if err != nil {
		fmt.Println("There was an error while adding the notification seen flag")
		fmt.Println(err)
	}
}

func pushNotificationForUser(NotificaionObject models.Notification, suffixTxt string) {
	// TODO: add the notification to db... having a prob heree..
	fmt.Println("Here in the NotificaionObject push func.")

	if(!(NotificaionObject.Actorid == NotificaionObject.Uuid)) {
		Client, ok := models.ClientPool.GetClient(NotificaionObject.Uuid);

		if ok {
			NotificaionObject.User_ = GetUserById(NotificaionObject.Actorid);
			NotificaionObject.Text = NotificaionObject.User_.UserName + " " + suffixTxt;

			var resp models.SocketMessage = NotificaionObject.EncodeToSocketResponse();

			Client.SendJSON(resp);
			
		} else {
			fmt.Println("Client offline.");
		}

		AddNewNotification(NotificaionObject);
	}

}

// type socketResp struct {
// 	Code int `json:"code"`
// 	Data interface{} `json:"data"`
// }