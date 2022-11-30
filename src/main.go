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

func run() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.Default())	
	router.POST("/login", login)
	router.GET("/getUserPosts", getUserPostsRoute)
	router.GET("/GetAllPosts", GetAllPostsRoute)
	router.GET("/query", getUsersRoute)
	router.GET("/:uuid", getUserByIdRoute)
	router.POST("/signup", signUp)
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

