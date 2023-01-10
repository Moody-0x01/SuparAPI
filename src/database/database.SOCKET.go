package database

import (
	"fmt"
	// "io"
	"github.com/Moody0101-X/Go_Api/models"
)

func HandleClientConnection(c *models.Client) {
	for {
		
		var sockmsg models.SocketMessage
		err := c.Conn.ReadJSON(&sockmsg)

		if err != nil {
			fmt.Println(err);
			break
		}

		switch sockmsg.Action {
			
			case models.NOTIFICATION:
				// var seenFlag models.Notification = sockmsg.Data.(models.Notification)
				var Id_ = int(sockmsg.Data.(map[string]interface{})["id"].(float64))
				SetSeenForNotification(Id_)
				break

			default:
				break

		}
	}
}