package main

import (
    // "net/http"
	// "github.com/gin-gonic/contrib/static"

    "fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

var (
	port string = ":8888"
)



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
	
	// Get routes.
	
	
	router.GET("/v2/getUserPosts", getUserPostsRoute) // gettting user post by id
	router.GET("/v2/GetAllPosts", GetAllPostsRoute) // getting all posts
	router.GET("/v2/query", getUsersRoute) // user look up by name
	router.GET("/v2/:uuid", getUserByIdRoute) // get user by id
	router.GET("/v2/getComments/:pid", getPostComments)
	router.GET("/v2/getLikes/:pid", getPostLikes)
	// router.Static("/", "./public")
	router.GET("/", index)

    router.NoRoute(index)

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

