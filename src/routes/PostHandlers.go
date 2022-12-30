package routes;

import (
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/Moody0101-X/Go_Api/models"
	"github.com/Moody0101-X/Go_Api/database"
	"github.com/Moody0101-X/Go_Api/crypto"
)

func Login(c *gin.Context) {
	var LoginForm models.UserLogin;
	
	c.BindJSON(&LoginForm);
	
	var resp models.Response

	if len(LoginForm.Token) > 0 {
		resp = database.AuthenticateUserJWT(LoginForm.Token)
	} else {
		if len(LoginForm.Password) > 0 && len(LoginForm.Email) > 0 {
			user, err := database.AuthenticateUserByEmailAndPwd(LoginForm.Password, LoginForm.Email)
			
			if err.Ok {
				resp = models.MakeServerResponse(200, user)
			} else {

				resp = models.MakeServerResponse(500, err.Text)
				fmt.Println("", resp.Data)
			}

		} else {
			resp = models.MakeServerResponse(500, "Missing request attributes, Email or password not specified.")
		}
	}

	// fmt.Println(resp)
	c.JSON(http.StatusOK, resp);
}

func SignUp(c *gin.Context) {
	
	var newUser models.User	
	c.BindJSON(&newUser);
	if isEmpty(newUser.Email) || isEmpty(newUser.PasswordHash) || isEmpty(newUser.UserName) {
		c.JSON(http.StatusOK, models.MakeServerResponse(500, "The server could not get the Email, password or user name. please check your request then try again L86"))
	} else {
		newUser.SetDefaults();
		// Hash the password.
		newUser.PasswordHash = crypto.Sha256_(newUser.PasswordHash)
		var Resp models.Response = database.AddUser(newUser) // Creates the user and sets the Token.
		c.JSON(http.StatusOK, Resp)
	}
}

func Update(c *gin.Context) {
	// This function can update. -bg, -bio, -img, -username
	// Token 		 string `json:"token"`
	// Img 		 string `json:"img"`
	// Bg 			 string `json:"bg"`
	// Bio 		 string `json:"bio"`
	// Address		 string `json:"addr"`

	var Data models.User
	c.BindJSON(&Data)
	if len(Data.Token) > 0 {
		
		AccessToken, Ok := crypto.GetTokenFromJwt(Data.Token)
		
		if Ok {
			var Ok bool = true;

			if !isEmpty(Data.Img) {
				e := database.UpdateUser("IMG", Data.Img, AccessToken)
				Ok = e.Ok
			}

			if !isEmpty(Data.Bio) {
				e := database.UpdateUser("BIO", Data.Bio, AccessToken)
				Ok = e.Ok
			}

			if !isEmpty(Data.Address){
				e := database.UpdateUser("ADDR", Data.Address, AccessToken)
				Ok = e.Ok
			}

			if !isEmpty(Data.Bg) {
				e := database.UpdateUser("BG", Data.Bg, AccessToken)
				Ok = e.Ok
			}

			if !isEmpty(Data.UserName) {
				e := database.UpdateUser("USERNAME", Data.UserName, AccessToken)
				Ok = e.Ok
			}

			if Ok {
				c.JSON(http.StatusOK, models.MakeServerResponse(200, "updated!"))
				return
			} else {
				c.JSON(http.StatusOK, models.MakeServerResponse(500, "Something went wrong.! 143"))
				return
			}

		} else {
			// return error. invalid token.
			c.JSON(http.StatusOK, models.MakeServerResponse(500, "The Token you have specified is invalid, try with other."))
			return
		}
	}
	
	c.JSON(http.StatusOK, models.MakeServerResponse(500, "No token was provided in the request. try agin with a token."))
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

	var post models.TokenizedPost
	c.BindJSON(&post)
	
	if isEmpty(post.Token) {
		c.JSON(http.StatusOK, models.MakeServerResponse(500, "no Token provided, try providing your secure token."))
		return
	}

	if post.Uuid == 0 {
		c.JSON(http.StatusOK, models.MakeServerResponse(500, "uuid field not present: Provide uuid"))
		return
	}

	if isEmpty(post.Text) && isEmpty(post.Img) {
		c.JSON(http.StatusOK, models.MakeServerResponse(500, "img and text are empty, provide some text for the post or an img"))
		return
	}

	_, Ok := crypto.GetTokenFromJwt(post.Token)
	


	if Ok {
		err := database.AddPost(post.Text, post.Img, post.Uuid)

		if err.Ok {
			c.JSON(http.StatusOK, models.MakeServerResponse(200, "success, added."))
			return
		}

		c.JSON(http.StatusOK, models.MakeServerResponse(500, "the user was not added. problem in db."))
		return
	}

	c.JSON(http.StatusOK, models.MakeServerResponse(500, "invalid access token. try with a valid token."))
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

	var post models.TokenizedPost
	c.BindJSON(&post)
	// - PostID
	// - Token
	// - Uuid

	if isEmpty(post.Token) {
		c.JSON(http.StatusOK, models.MakeServerResponse(500, "A token is required {token: v}"))
		return
	}

	if post.Uuid == 0  || post.PostID == 0 {
		c.JSON(http.StatusOK, models.MakeServerResponse(500, "a required argument is messing, check the uuid and id_"))
		return
	}

	token, Ok := crypto.GetTokenFromJwt(post.Token)
	var resp models.Response
	
	if Ok {
		resp = database.DeleteUserPost(post.PostID, post.Uuid, token)
		c.JSON(http.StatusOK, resp)
		return
	}

	c.JSON(http.StatusOK, models.MakeServerResponse(500, "invalid access token. try with a valid token."))
}

func AddCommentRoute(c *gin.Context) {
	/*  Excpects: 
		{
			"post_id": str
			"token": v, // SO a user only can add his own comments.. SECURITY LESS GOO
			"uuid": int,
			"text": str,
		}  
	*/

	var CommentRoutePostedData models.TokenizedComment
	c.BindJSON(&CommentRoutePostedData)

	if isEmpty(CommentRoutePostedData.Token) || CommentRoutePostedData.Uuid == 0 || CommentRoutePostedData.Post_id == 0 || isEmpty(CommentRoutePostedData.Text) {
		c.JSON(http.StatusOK, models.MakeServerResponse(400, "Bad request, token | post_id | text | uuid is Missing"))
		return
	}

	AccessToken, Ok := crypto.GetTokenFromJwt(CommentRoutePostedData.Token)

	if Ok {
		// passing the other data to add the comment.
		result := database.Add_comment(CommentRoutePostedData.Uuid, CommentRoutePostedData.Text, CommentRoutePostedData.Post_id, AccessToken, CommentRoutePostedData.Post_owner_id)
		
		if result.Ok {
			c.JSON(http.StatusOK, models.MakeServerResponse(200, result.Text))
			return
		}
		
		c.JSON(http.StatusOK, models.MakeServerResponse(500, result.Text))
		return
	}

	c.JSON(http.StatusOK, models.MakeServerResponse(401, "The token sent is not valid!"))
}

func AddLikeRoute(c *gin.Context) {
	/*  Excpects:
		{
			"post_id": str
			"token": v, // SO a user only can add his own comments.. SECURITY LESS GOO
			"uuid": int
		}  
	}  */

	var LikeRoutePostedData models.TokenizedLike;
	c.BindJSON(&LikeRoutePostedData)

	if isEmpty(LikeRoutePostedData.Token) || LikeRoutePostedData.Uuid == 0 ||  LikeRoutePostedData.Post_id == 0 {
		c.JSON(http.StatusOK, models.MakeServerResponse(400, "Bad request, token | post_id | uuid is Missing"))
		return 
	}

	AccessToken, Ok := crypto.GetTokenFromJwt(LikeRoutePostedData.Token)

	if Ok {
		// passing the other data to add the Like.
		fmt.Println("uuid that was parsed: ", LikeRoutePostedData.Post_owner_id);

		result := database.Add_like(LikeRoutePostedData.Uuid, LikeRoutePostedData.Post_id, AccessToken, LikeRoutePostedData.Post_owner_id)
		
		if result.Ok {
			var data models.Like;
			
			data.Post_id = LikeRoutePostedData.Post_id
			data.Uuid = LikeRoutePostedData.Uuid
			data.User_ = database.GetUserById(data.Uuid)

			c.JSON(http.StatusOK, models.MakeServerResponse(200, data))
			return
		}
		
		c.JSON(http.StatusOK, models.MakeServerResponse(500, result.Text))
		return
	}

	c.JSON(http.StatusOK, models.MakeServerResponse(401, "The token sent is not valid!"))
}

func RemoveLikeRoute(c *gin.Context) {

	var LikeRoutePostedData models.TokenizedLike;
	c.BindJSON(&LikeRoutePostedData)

	if isEmpty(LikeRoutePostedData.Token) || LikeRoutePostedData.Uuid == 0 ||  LikeRoutePostedData.Post_id == 0 {
		c.JSON(http.StatusOK, models.MakeServerResponse(400, "Bad request, token | post_id | uuid is Missing"))
		return 
	}

	AccessToken, Ok := crypto.GetTokenFromJwt(LikeRoutePostedData.Token)

	if Ok {
		// passing the other data to add the Like.
		result := database.Remove_like(LikeRoutePostedData.Uuid, LikeRoutePostedData.Post_id, AccessToken)
		
		if result.Ok {
			c.JSON(http.StatusOK, models.MakeServerResponse(200, result.Text))
			return
		}
		
		c.JSON(http.StatusOK, models.MakeServerResponse(500, result.Text))
		return
	}

	c.JSON(http.StatusOK, models.MakeServerResponse(401, "The token sent is not valid!"))

}

func FollowRoute(c *gin.Context) {
	//  Gets the follower_id and the one that wants to be added.
	var UData models.TFollow;

	c.BindJSON(&UData)

	if isEmpty(UData.UToken) || UData.Follower_id == 0 || UData.Followed_id == 0 {
		c.JSON(http.StatusOK, models.MakeServerResponse(400, "Bad request, token | follower_id | followed_id is Missing"))
		return 
	}

	// type  struct {
	// 	Follower_id		int `json:"follower_id"`
	// 	Followed_id		int `json:"followed_id"`
	// 	UToken			string `json:"token"`
	// }

	AccessToken, Ok := crypto.GetTokenFromJwt(UData.UToken)

	if Ok {
		// passing the other data to add the Like.
		result := database.Follow(UData.Follower_id, UData.Followed_id, AccessToken)
		
		if result.Ok {
			c.JSON(http.StatusOK, models.MakeServerResponse(200, result.Text))
			return
		}
		
		c.JSON(http.StatusOK, models.MakeServerResponse(500, result.Text))
		return
	}

	c.JSON(http.StatusOK, models.MakeServerResponse(401, "The token sent is not valid!"))

	// notImplemented(c);
}

func UnfollowRoute(c *gin.Context) {
	
	var UData models.TFollow;
	c.BindJSON(&UData)

	if isEmpty(UData.UToken) || UData.Follower_id == 0 || UData.Followed_id == 0 {
		c.JSON(http.StatusOK, models.MakeServerResponse(400, "Bad request, token | follower_id | followed_id is Missing"))
		return 
	}

	// type  struct {
	// 	Follower_id		int `json:"follower_id"`
	// 	Followed_id		int `json:"follower_id"`
	// 	UToken			string `json:"token"`
	// }

	AccessToken, Ok := crypto.GetTokenFromJwt(UData.UToken)

	if Ok {
		// passing the other data to add the unfollow event..
		result := database.Unfollow(UData.Follower_id, UData.Followed_id, AccessToken)

		if result.Ok {
			c.JSON(http.StatusOK, models.MakeServerResponse(200, result.Text))
			return
		}
		
		c.JSON(http.StatusOK, models.MakeServerResponse(500, result.Text))
		
		return
	}

	c.JSON(http.StatusOK, models.MakeServerResponse(401, "The token sent is not valid!"))
	// notImplemented(c);
}