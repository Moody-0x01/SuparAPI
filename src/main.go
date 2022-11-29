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

func run() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.Default())
	

	
	router.POST("/login", login)
	router.POST("/Test", LOGIN)
	router.GET("/getUserPosts", getUserPostsRoute)
	router.GET("/GetAllPosts", GetAllPostsRoute)
	router.GET("/query", getUsersRoute)
	router.GET("/:uuid", getUserByIdRoute)

	fmt.Println("Serving in port 8888")
	router.Run(":8888")
}


func main() {

	err := initializeDb();

	if err != nil {
        fmt.Println("Error opening the database! ", err.Error())
        return
    }
	
	run()
}

