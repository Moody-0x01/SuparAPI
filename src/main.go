package main

import (
    // "net/http"
	// "github.com/gin-gonic/contrib/static"

    "fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"

)

var (
	port string = ":8888"
	socketClients []client
)

type client struct {
	Addr string
	Uuid int
	Conn *websocket.Conn
	IsOpen bool
}

func newClient(Addr string, Uuid int, Conn *websocket.Conn, IsOpen bool) {
	var New client
	New.Addr = Addr
	New.Uuid = Uuid
	New.Conn = Conn
	New.IsOpen = IsOpen
	socketClients = append(socketClients, New)
	go New.handleConn()
}

func (c *client) logClient() {
	fmt.Println("addr: ", c.Addr)
	fmt.Println("uuid: ", c.Uuid)
}

var upgrader = websocket.Upgrader{
    //check origin will check the cross region source (note : please not using in production)
    ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
        //Here we just allow the chrome extension client accessable (you should check this verify accourding your client source)
		return true
	},
}



func (c *client) handleConn() {
	
	for {	
		mt, message, err := c.Conn.ReadMessage()

		if err != nil {
			fmt.Println(err)
			break
		}
			
		defer c.Conn.Close()

		for i := 0; i < len(socketClients); i++ {
			
			if socketClients[i] != *c {
				var conn = socketClients[i].Conn
				
				NewMsg := []byte(strconv.Itoa(socketClients[i].Uuid) + "said: " + string(message))
					
				err = conn.WriteMessage(mt, NewMsg)
				
				if err != nil {
					fmt.Println(err)
					break
				}
			}
		}
	}
}

func WebSocketRoute(c *gin.Context) {
	//upgrade get request to websocket protocol
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
	} else {
		fmt.Println("new uuid: ", uuid_)
	}

	newClient(ws.RemoteAddr().String(), uuid_, ws, true)
}


func RequestCancelRecover() gin.HandlerFunc {
	
	return func(c *gin.Context) {
		defer func() {
			
			if err := recover(); err != nil {
				fmt.Println("client cancel the request")
				c.Request.Context().Done()
			}

		}()
		
		c.Next()
	}

}


func run() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(cors.Default())	
	router.Use(gin.Logger(), RequestCancelRecover())
	// HTML/JS/CSS/IMG loaders
	router.Static("/static", "./public/static")
	router.Static("/img", "./public/img")	
	router.LoadHTMLGlob("public/*.html")

	// POST routes.
	router.POST("/v2/login", login) // login and get a token for the updating/creation/deletion of personal data.
	router.POST("/v2/update", update) // Updating user's information by token
	router.POST("/v2/NewPost", NewPost) // adding a post by token.
	router.POST("/v2/DeletePost", DeletePost) // Deleting a post by token
	router.POST("/v2/signup", signUp) // Making new account
	router.POST("/v2/comment", addCommentRoute) // For likes
	router.POST("/v2/like", addLikeRoute) // For comments
	router.POST("/v2/like/remove", RemoveLikeRoute)
	router.POST("/v2/follow", followRoute)
	router.POST("/v2/unfollow", unfollowRoute)
	
	// Get routes.
	router.GET("/v2/getUserPosts", getUserPostsRoute) // gettting user post by id
	router.GET("/v2/GetAllPosts", GetAllPostsRoute) // getting all posts
	router.GET("/v2/query", getUsersRoute) // user look up by name
	router.GET("/v2/:uuid", getUserByIdRoute) // get user by id
	router.GET("/v2/getFollowers/:uuid", getUserFollowersById)
	router.GET("/v2/getComments/:pid", getPostComments)
	router.GET("/v2/getLikes/:pid", getPostLikes)
	// router.Static("/", "./public")
	router.GET("/", index)
    router.NoRoute(index)

    // Socket routes.
    router.GET("/v2/ws", WebSocketRoute)
	// running the server.
	fmt.Println("Serving in port", port)
	
	router.Run(port)
}

func main() {
	err, path := initializeDb();
	
	if err != nil {
        fmt.Println("Error opening the database! ", err.Error())
        return
    }

    fmt.Println("connected to db: ", path)

	run()
}

