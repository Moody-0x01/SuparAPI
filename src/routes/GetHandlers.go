package routes;

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
	"github.com/Moody0101-X/Go_Api/models"
	"github.com/Moody0101-X/Go_Api/database"
)

func GetAllPostsRoute(c *gin.Context) {
	All := database.GetAllPosts()
	var res models.Response = models.MakeServerResponse(200, All)
	c.JSON(http.StatusOK, res)
}

func GetUserPostsRoute(c *gin.Context) {
	var id string = GetFieldFromContext(c, "id_")
	id_, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(500, models.MakeServerResponse(0, "Server error"))
		return 
	}

	var UserPosts []models.Post = database.GetUserPostById(id_)
	c.JSON(http.StatusOK, models.MakeServerResponse(200, UserPosts))

}

func GetUsersRoute(c *gin.Context) {
	var q string = GetFieldFromContext(c, "q")
	var uuid string = GetFieldFromContext(c, "uuid")
	var uuid_ int;
	var flag bool = false;
	if !isEmpty(uuid) {
		flag = true;
		temp, err := strconv.Atoi(uuid)
		uuid_ = temp

		if err != nil {
			c.JSON(http.StatusOK, models.MakeServerResponse(400, "make sure that uuid is a number."))
			return
		}
	}

	var Users []models.AUser;

	if q != "" {
		if flag {
			Users = database.GetUsersByQuery(q, uuid_)
			c.JSON(http.StatusOK, models.MakeServerResponse(200, Users))
			return
		} else {
			Users = database.GetUsersByQuery(q, flag)	
			c.JSON(http.StatusOK, models.MakeServerResponse(200, Users))
			return
		}
	} else {
		if flag {
			Users = database.GetUsers(uuid_)
			c.JSON(http.StatusOK, models.MakeServerResponse(200, Users))
			return
		} else {
			Users = database.GetUsers(flag)
			c.JSON(http.StatusOK, models.MakeServerResponse(200, Users))
			return
		}
	}
}

func GetUserByIdRoute(c *gin.Context) {
	
	var uuid string = GetFieldFromContext(c, "uuid")
	var user_id string = GetFieldFromContext(c, "user")

	uuid_, err := strconv.Atoi(uuid)

	if err != nil {
		c.JSON(500, models.MakeServerResponse(400, "uuid should be a number"))
		return 
	}

	var User models.AUser = database.GetUserById(uuid_)
	
	if(!isEmpty(user_id)) {
		user_ID, err := strconv.Atoi(user_id)

		if err != nil {
			c.JSON(500, models.MakeServerResponse(400, "user get parameter should be a number for this request to successed"))
			return 
		}

		User.IsFollowed = database.IsFollowing(User.Id_, user_ID)
	}

	var res models.Response = models.MakeServerResponse(200, User)
	
	c.JSON(http.StatusOK, res)
}

func GetPostComments(c *gin.Context) {
	var pid string = c.Param("pid")
	
	pid_, err := strconv.Atoi(pid)

	if err != nil {
		c.JSON(http.StatusOK, models.MakeServerResponse(400, "bad request, make sure post_id is an integer"))
		return 
	}
	
	var comments []models.Comment = database.Get_comments(pid_);

	c.JSON(http.StatusOK, models.MakeServerResponse(200, comments))
}

func GetPostLikes(c *gin.Context) {
	var pid string = c.Param("pid")
	
	pid_, err := strconv.Atoi(pid)

	if err != nil {
		c.JSON(http.StatusOK, models.MakeServerResponse(400, "bad request, make sure post_id is an integer"))
		return 
	}
	
	var likes []models.Like = database.Get_likes(pid_);

	c.JSON(http.StatusOK, models.MakeServerResponse(200, likes))
}


func GetUserFollowingsById(c *gin.Context) {
	
	var uuid string = c.Param("uuid")
	
	uuid_, err := strconv.Atoi( uuid )

	if err != nil {
		c.JSON(http.StatusOK, models.MakeServerResponse(400, "bad request, make sure post_id is an integer"))
		return
	}

	
	var followers []int = database.GetFollowings(uuid_)
	
	c.JSON(http.StatusOK, models.MakeServerResponse(200, followers))
}

func GetUserFollowersById(c *gin.Context) {
	
	var uuid string = c.Param("uuid")
	
	uuid_, err := strconv.Atoi( uuid )

	if err != nil {
		c.JSON(http.StatusOK, models.MakeServerResponse(400, "bad request, make sure post_id is an integer"))
		return
	}

	var followers []int = database.GetFollowers(uuid_)
	
	c.JSON(http.StatusOK, models.MakeServerResponse(200, followers))
}

func GetPostByPostidRoute(c *gin.Context) {
	var pid string = c.Param("pid")
	
	pid_, err := strconv.Atoi(pid)

	if err != nil {
		c.JSON(http.StatusOK, models.MakeServerResponse(400, "bad request, make sure post_id is an pid"))
		return
	}
	
	var Post_ models.Post = database.GetPostById(pid_)
	
	c.JSON(http.StatusOK, models.MakeServerResponse(200, Post_))
}

func GetAllNotificationsRoute(c *gin.Context) {

	var uuid string = c.Param("uuid")
	uuid_, err := strconv.Atoi( uuid )

	if err != nil {
		c.JSON(http.StatusOK, models.MakeServerResponse(400, "bad request, make sure post_id is an integer"))
		return
	}

	var Notifications []models.Notification = database.GetAllNotifications(uuid_)

	c.JSON(http.StatusOK, models.MakeServerResponse(200, Notifications))
}