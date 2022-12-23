package main;

import (
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func GetAllPostsRoute(c *gin.Context) {
	All := GetAllPosts()
	var res Response = MakeServerResponse(200, All)
	c.JSON(http.StatusOK, res)
}
	

func getUserPostsRoute(c *gin.Context) {
	var id string = GetFieldFromContext(c, "id_")
	id_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(500, MakeServerResponse(0, "Server error"))
		return 
	}

	var UserPosts []Post = getUserPostById(id_)
	c.JSON(http.StatusOK, MakeServerResponse(200, UserPosts))
}


func getUsersRoute(c *gin.Context) {
	var q string = GetFieldFromContext(c, "q")
	var uuid string = GetFieldFromContext(c, "uuid")
	var uuid_ int;
	var flag bool = false;
	if !isEmpty(uuid) {
		flag = true;
		temp, err := strconv.Atoi(uuid)
		uuid_ = temp

		if err != nil {
			c.JSON(http.StatusOK, MakeServerResponse(400, "make sure that uuid is a number."))
			return
		}
	}

	var Users []AUser;

	if q != "" {
		if flag {
			Users = getUsersByQuery(q, uuid_)
			c.JSON(http.StatusOK, MakeServerResponse(200, Users))
			return
		} else {
			Users = getUsersByQuery(q, flag)	
			c.JSON(http.StatusOK, MakeServerResponse(200, Users))
			return
		}
	} else {
		if flag {
			Users = getUsers(uuid_)
			c.JSON(http.StatusOK, MakeServerResponse(200, Users))
			return
		} else {
			Users = getUsers(flag)
			c.JSON(http.StatusOK, MakeServerResponse(200, Users))
			return
		}
	}
}

func getUserByIdRoute(c *gin.Context) {
	
	var uuid string = GetFieldFromContext(c, "uuid")
	var user_id string = GetFieldFromContext(c, "user")

	uuid_, err := strconv.Atoi(uuid)

	if err != nil {
		c.JSON(500, MakeServerResponse(400, "uuid should be a number"))
		return 
	}

	var User AUser = getUserById(uuid_)
	
	if(!isEmpty(user_id)) {
		user_ID, err := strconv.Atoi(user_id)

		if err != nil {
			c.JSON(500, MakeServerResponse(400, "user get parameter should be a number for this request to successed"))
			return 
		}

		User.IsFollowed = isFollowing(User.Id_, user_ID)
	}

	var res Response = MakeServerResponse(200, User)
	
	c.JSON(http.StatusOK, res)
}

/* AUTHENTICATION AND OPERATIONS */

func login(c *gin.Context) {
	var LoginForm UserLogin;
	
	c.BindJSON(&LoginForm);
	
	var resp Response

	if len(LoginForm.Token) > 0 {
		resp = AuthenticateUserJWT(LoginForm.Token)
	} else {
		if len(LoginForm.Password) > 0 && len(LoginForm.Email) > 0 {
			user, err := AuthenticateUserByEmailAndPwd(LoginForm.Password, LoginForm.Email)
			
			if err.Ok {
				resp = MakeServerResponse(200, user)
			} else {

				resp = MakeServerResponse(500, err.Text)
				fmt.Println("", resp.Data)
			}

		} else {
			resp = MakeServerResponse(500, "Missing request attributes, Email or password not specified.")
		}
	}

	// fmt.Println(resp)
	c.JSON(http.StatusOK, resp);
}


func signUp(c *gin.Context) {
	
	var newUser User	
	c.BindJSON(&newUser);
	if isEmpty(newUser.Email) || isEmpty(newUser.PasswordHash) || isEmpty(newUser.UserName) {
		c.JSON(http.StatusOK, MakeServerResponse(500, "The server could not get the Email, password or user name. please check your request then try again L86"))
	} else {
		newUser.setDefaults();
		// Hash the password.
		newUser.PasswordHash = sha256_(newUser.PasswordHash)
		var Resp Response = AddUser(newUser) // Creates the user and sets the Token.
		c.JSON(http.StatusOK, Resp)
	}
}

/*------------------------------------------------------------------------------------------------------------------------*/

//TODO Make AddPost function.
//TODO Make update generic function.

func update(c *gin.Context) {
	// This function can update. -bg, -bio, -img, -username
	// Token 		 string `json:"token"`
	// Img 		 string `json:"img"`
	// Bg 			 string `json:"bg"`
	// Bio 		 string `json:"bio"`
	// Address		 string `json:"addr"`

	var Data User
	c.BindJSON(&Data)
	if len(Data.Token) > 0 {
		AccessToken, Ok := GetTokenFromJwt(Data.Token)
		
		if Ok {
			var Ok bool = true;

			if !isEmpty(Data.Img) {
				e := updateUser("IMG", Data.Img, AccessToken)
				Ok = e.Ok
			}

			if !isEmpty(Data.Bio) {
				e := updateUser("BIO", Data.Bio, AccessToken)
				Ok = e.Ok
			}

			if !isEmpty(Data.Address){
				e := updateUser("ADDR", Data.Address, AccessToken)
				Ok = e.Ok
			}

			if !isEmpty(Data.Bg) {
				e := updateUser("BG", Data.Bg, AccessToken)
				Ok = e.Ok
			}

			if !isEmpty(Data.UserName) {
				e := updateUser("USERNAME", Data.UserName, AccessToken)
				Ok = e.Ok
			}

			if Ok {
				c.JSON(http.StatusOK, MakeServerResponse(200, "updated!"))
				return
			} else {
				c.JSON(http.StatusOK, MakeServerResponse(500, "Something went wrong.! 143"))
				return
			}

		} else {
			// return error. invalid token.
			c.JSON(http.StatusOK, MakeServerResponse(500, "The Token you have specified is invalid, try with other."))
			return
		}
	}
	
	c.JSON(http.StatusOK, MakeServerResponse(500, "No token was provided in the request. try agin with a token."))
}

func NewPost(c *gin.Context) {
	/*
	Expectation:
		json = {
			"Token": v,
			"uuid": v,
			"text": v,
			"img": v
		}
	*/
	// This function creates a post for the user.
	// ID INTEGER PRIMARY KEY AUTOINCREMENT,
	// USER_ID INTEGER,
	// Text TEXT,
	// IMG TEXT,

	var post TokenizedPost
	c.BindJSON(&post)
	
	if isEmpty(post.Token) {
		c.JSON(http.StatusOK, MakeServerResponse(500, "no Token provided, try providing your secure token."))
		return
	}

	if post.Uuid == 0 {
		c.JSON(http.StatusOK, MakeServerResponse(500, "uuid field not present: Provide uuid"))
		return
	}

	if isEmpty(post.Text) && isEmpty(post.Img) {
		c.JSON(http.StatusOK, MakeServerResponse(500, "img and text are empty, provide some text for the post or an img"))
		return
	}

	_, Ok := GetTokenFromJwt(post.Token)
	


	if Ok {
		err := AddPost(post.Text, post.Img, post.Uuid)

		if err.Ok {
			c.JSON(http.StatusOK, MakeServerResponse(200, "success, added."))
			return
		}

		c.JSON(http.StatusOK, MakeServerResponse(500, "the user was not added. problem in db."))
		return
	}

	c.JSON(http.StatusOK, MakeServerResponse(500, "invalid access token. try with a valid token."))
}

func DeletePost(c *gin.Context) {
	/*
	
	expecting: 
		Json {
			"id_"
			"token"
			"uuid"
		}

	*/

	var post TokenizedPost
	c.BindJSON(&post)
	// - PostID
	// - Token
	// - Uuid

	if isEmpty(post.Token) {
		c.JSON(http.StatusOK, MakeServerResponse(500, "A token is required {token: v}"))
		return
	}

	if post.Uuid == 0  || post.PostID == 0 {
		c.JSON(http.StatusOK, MakeServerResponse(500, "a required argument is messing, check the uuid and id_"))
		return
	}

	token, Ok := GetTokenFromJwt(post.Token)
	var resp Response
	
	if Ok {
		resp = DeleteUserPost(post.PostID, post.Uuid, token)
		c.JSON(http.StatusOK, resp)
		return
	}

	c.JSON(http.StatusOK, MakeServerResponse(500, "invalid access token. try with a valid token."))
}


func getPostComments(c *gin.Context) {
	var pid string = c.Param("pid")
	
	pid_, err := strconv.Atoi(pid)

	if err != nil {
		c.JSON(http.StatusOK, MakeServerResponse(400, "bad request, make sure post_id is an integer"))
		return 
	}
	
	var comments []Comment = get_comments(pid_);

	c.JSON(http.StatusOK, MakeServerResponse(200, comments))
}

func getPostLikes(c *gin.Context) {
	var pid string = c.Param("pid")
	
	pid_, err := strconv.Atoi(pid)

	if err != nil {
		c.JSON(http.StatusOK, MakeServerResponse(400, "bad request, make sure post_id is an integer"))
		return 
	}
	
	var likes []Like = get_likes(pid_);

	c.JSON(http.StatusOK, MakeServerResponse(200, likes))
}

// type TokenizedComment struct {
// 	Post_id		 int `json:"post_id"`
// 	Uuid		 int `json:"uuid"`
// 	Text		 string `json:"text"`
// 	Token        string `json:"token"`
// }

// type TokenizedLike struct {
// 	Post_id		 int `json:"post_id"`
// 	Uuid		 int `json:"uuid"`
// 	Token        string `json:"token"`
// }

func addCommentRoute(c *gin.Context) {
	/*  Excpects: 
		{
			"post_id": str
			"token": v, // SO a user only can add his own comments.. SECURITY LESS GOO
			"uuid": int,
			"text": str,
		}  
	*/

	var CommentRoutePostedData TokenizedComment
	c.BindJSON(&CommentRoutePostedData)

	if isEmpty(CommentRoutePostedData.Token) || CommentRoutePostedData.Uuid == 0 || CommentRoutePostedData.Post_id == 0 || isEmpty(CommentRoutePostedData.Text) {
		c.JSON(http.StatusOK, MakeServerResponse(400, "Bad request, token | post_id | text | uuid is Missing"))
		return
	}

	AccessToken, Ok := GetTokenFromJwt(CommentRoutePostedData.Token)

	if Ok {
		// passing the other data to add the comment.
		result := add_comment(CommentRoutePostedData.Uuid, CommentRoutePostedData.Text, CommentRoutePostedData.Post_id, AccessToken)
		if result.Ok {
			c.JSON(http.StatusOK, MakeServerResponse(200, result.Text))
			return
		}
		
		c.JSON(http.StatusOK, MakeServerResponse(500, result.Text))
		return
	}

	c.JSON(http.StatusOK, MakeServerResponse(401, "The token sent is not valid!"))
}

func addLikeRoute(c *gin.Context) {
	/*  Excpects:
		{
			"post_id": str
			"token": v, // SO a user only can add his own comments.. SECURITY LESS GOO
			"uuid": int
		}  
	}  */

	var LikeRoutePostedData TokenizedLike;
	c.BindJSON(&LikeRoutePostedData)

	if isEmpty(LikeRoutePostedData.Token) || LikeRoutePostedData.Uuid == 0 ||  LikeRoutePostedData.Post_id == 0 {
		c.JSON(http.StatusOK, MakeServerResponse(400, "Bad request, token | post_id | uuid is Missing"))
		return 
	}

	AccessToken, Ok := GetTokenFromJwt(LikeRoutePostedData.Token)

	if Ok {
		// passing the other data to add the Like.
		result := add_like(LikeRoutePostedData.Uuid, LikeRoutePostedData.Post_id, AccessToken)
		
		if result.Ok {
			var data Like;
			
			data.Post_id = LikeRoutePostedData.Post_id
			data.Uuid = LikeRoutePostedData.Uuid
			data.User_ = getUserById(data.Uuid)

			c.JSON(http.StatusOK, MakeServerResponse(200, data))
			return
		}
		
		c.JSON(http.StatusOK, MakeServerResponse(500, result.Text))
		return
	}

	c.JSON(http.StatusOK, MakeServerResponse(401, "The token sent is not valid!"))
}

func RemoveLikeRoute(c *gin.Context) {

	var LikeRoutePostedData TokenizedLike;
	c.BindJSON(&LikeRoutePostedData)

	if isEmpty(LikeRoutePostedData.Token) || LikeRoutePostedData.Uuid == 0 ||  LikeRoutePostedData.Post_id == 0 {
		c.JSON(http.StatusOK, MakeServerResponse(400, "Bad request, token | post_id | uuid is Missing"))
		return 
	}

	AccessToken, Ok := GetTokenFromJwt(LikeRoutePostedData.Token)

	if Ok {
		// passing the other data to add the Like.
		result := remove_like(LikeRoutePostedData.Uuid, LikeRoutePostedData.Post_id, AccessToken)
		
		if result.Ok {
			c.JSON(http.StatusOK, MakeServerResponse(200, result.Text))
			return
		}
		
		c.JSON(http.StatusOK, MakeServerResponse(500, result.Text))
		return
	}

	c.JSON(http.StatusOK, MakeServerResponse(401, "The token sent is not valid!"))

}


func getUserFollowingsById(c *gin.Context) {
	
	var uuid string = c.Param("uuid")
	
	uuid_, err := strconv.Atoi( uuid )

	if err != nil {
		c.JSON(http.StatusOK, MakeServerResponse(400, "bad request, make sure post_id is an integer"))
		return
	}

	
	var followers []int = getFollowings(uuid_)
	
	c.JSON(http.StatusOK, MakeServerResponse(200, followers))
}

func getUserFollowersById(c *gin.Context) {
	
	var uuid string = c.Param("uuid")
	
	uuid_, err := strconv.Atoi( uuid )

	if err != nil {
		c.JSON(http.StatusOK, MakeServerResponse(400, "bad request, make sure post_id is an integer"))
		return
	}

	var followers []int = getFollowers(uuid_)
	
	c.JSON(http.StatusOK, MakeServerResponse(200, followers))
}


func followRoute(c *gin.Context) {
	//  Gets the follower_id and the one that wants to be added.
	var UData TFollow;

	c.BindJSON(&UData)

	if isEmpty(UData.UToken) || UData.Follower_id == 0 || UData.Followed_id == 0 {
		c.JSON(http.StatusOK, MakeServerResponse(400, "Bad request, token | follower_id | followed_id is Missing"))
		return 
	}

	// type  struct {
	// 	Follower_id		int `json:"follower_id"`
	// 	Followed_id		int `json:"followed_id"`
	// 	UToken			string `json:"token"`
	// }

	AccessToken, Ok := GetTokenFromJwt(UData.UToken)

	if Ok {
		// passing the other data to add the Like.
		result := follow(UData.Follower_id, UData.Followed_id, AccessToken)
		
		if result.Ok {
			c.JSON(http.StatusOK, MakeServerResponse(200, result.Text))
			return
		}
		
		c.JSON(http.StatusOK, MakeServerResponse(500, result.Text))
		return
	}

	c.JSON(http.StatusOK, MakeServerResponse(401, "The token sent is not valid!"))

	// notImplemented(c);
}

func unfollowRoute(c *gin.Context) {
	
	var UData TFollow;
	c.BindJSON(&UData)

	if isEmpty(UData.UToken) || UData.Follower_id == 0 || UData.Followed_id == 0 {
		c.JSON(http.StatusOK, MakeServerResponse(400, "Bad request, token | follower_id | followed_id is Missing"))
		return 
	}

	// type  struct {
	// 	Follower_id		int `json:"follower_id"`
	// 	Followed_id		int `json:"follower_id"`
	// 	UToken			string `json:"token"`
	// }

	AccessToken, Ok := GetTokenFromJwt(UData.UToken)

	if Ok {
		// passing the other data to add the unfollow event..
		result := unfollow(UData.Follower_id, UData.Followed_id, AccessToken)

		if result.Ok {
			c.JSON(http.StatusOK, MakeServerResponse(200, result.Text))
			return
		}
		
		c.JSON(http.StatusOK, MakeServerResponse(500, result.Text))
		
		return
	}

	c.JSON(http.StatusOK, MakeServerResponse(401, "The token sent is not valid!"))
	// notImplemented(c);
}

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusOK, MakeServerResponse(100, "Not implemented!"))	
}
