package database

import (
	"fmt"
	// "io"
	"encoding/json"
	"github.com/Moody0101-X/Go_Api/models"
)

func HandleClientConnection(c *models.Client) {
	// TODO: Add a new case in the switch that will handle the messages that will be routed to another u connexion.Running
Loop:
	for {
		
		var sockmsg models.SocketMessage
		err := c.Conn.ReadJSON(&sockmsg)

		if err != nil {
			fmt.Println(err);
			break Loop
		}

		switch sockmsg.Action {
			
			case models.NOTIFICATION:
				// var seenFlag models.Notification = sockmsg.Data.(models.Notification)
				var Id_ = int(sockmsg.Data.(map[string]interface{})["id"].(float64))
				SetSeenForNotification(Id_)
				break

			case models.MSG:
				// Deserialize the data of the sent req, then send it to the corresponding other. then add it to the db!
				var New models.UMessage;
				jsonData, _ := json.Marshal(sockmsg.Data.(map[string]interface{}));	
				json.Unmarshal(jsonData, &New);
				SendMessage(c, New);
				
				break
			
			default:
				break

		}
	}

}