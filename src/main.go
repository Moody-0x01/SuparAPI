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



// func RequestCancelRecover() gin.HandlerFunc {
	
// 	return func(c *gin.Context) {
// 		defer func() {
			
// 			if err := recover(); err != nil {
// 				fmt.Println("client cancel the request")
// 				c.Request.Context().Done()
// 			}

// 		}()
		
// 		c.Next()
// 	}

// }


func run() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(cors.Default())	
	// router.Use(gin.Logger(), RequestCancelRecover())
	
	// POST routes.
	router.POST("/login", login) // login and get a token for the updating/creation/deletion of personal data.
	router.POST("update", update) // Updating user's information by token
	router.POST("/NewPost", NewPost) // adding a post by token.
	router.POST("/DeletePost", DeletePost) // Deleting a post by token
	router.POST("/signup", signUp) // Making new account

	// Get routes.
	router.GET("/getUserPosts", getUserPostsRoute) // gettting user post by id
	router.GET("/GetAllPosts", GetAllPostsRoute) // getting all posts
	router.GET("/query", getUsersRoute) // user look up by name
	router.GET("/:uuid", getUserByIdRoute) // get user by id
	
	// running the server.
	fmt.Println("Serving in port", port) 	
	router.Run(port)
}

func main() {
	err := initializeDb();
	
	if err != nil {
        fmt.Println("Error opening the database! ", err.Error())
        return
    }

	run()
}

