package socketOperations

import (
	"net/http"
	"strconv"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/gin-gonic/gin"
	"github.com/Moody0101-X/Go_Api/models"
	"github.com/Moody0101-X/Go_Api/database"
)


var upgrader = websocket.Upgrader{
    //check origin will check the cross region source (note : please not using in production)
    ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
        //Here we just allow the chrome extension client accessable (you should check this verify accourding your client source)
		return true
	},
}

func NotificationServer(c *gin.Context) {
	
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	
	if err != nil {
		fmt.Println(err)
		return
	}

	_, uuid, err := ws.ReadMessage()

	uuid_, err := strconv.Atoi(string(uuid))
	
	if err != nil {
		fmt.Println("err: ", err)
		return 
	}

	client, ok := models.NewClient(ws.RemoteAddr().String(), uuid_, ws);
	if(ok) {
		go database.HandleConnextionForNotifications(client);	
	} else {
		fmt.Println("Could not register client with uuid: ", uuid_);
	}
}