package main

import (
    "fmt"
    // "net/http"
	// "github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
//	"encoding/json"
	// "strconv"
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
	// POST routes.
	router.POST("/login", login)
	router.POST("update", update)
	router.POST("/NewPost", NewPost)
	router.POST("/signup", signUp)

	// Get routes.
	router.GET("/getUserPosts", getUserPostsRoute)
	router.GET("/GetAllPosts", GetAllPostsRoute)
	router.GET("/query", getUsersRoute)
	router.GET("/:uuid", getUserByIdRoute)
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

