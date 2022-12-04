package main;

import (
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

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
	var Users []AUser;
	
	if q != "" {
		Users = getUsersByQuery(q)
	} else {
		Users = getUsers()
	}
	
	c.JSON(http.StatusOK, MakeServerResponse(200, Users))
}

func getUserByIdRoute(c *gin.Context) {
	
	var uuid string = c.Param("uuid")
	uuid_, err := strconv.Atoi(uuid)

	if err != nil {
		c.JSON(500, MakeServerResponse(0, "Server error"))
		return 
	}

	var User AUser = getUserById(uuid_)
	var res Response = MakeServerResponse(200, User)
	
	c.JSON(http.StatusOK, res)
}

/* AUTHENTICATION AND OPERATIONS */
/*
Implemented: Login, Sign Up.
Not implemented: data access (Update, add, delete)
*/
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
				e :=updateUser("BG", Data.Bg, AccessToken)
				Ok = e.Ok
			}

			if !isEmpty(Data.UserName) {
				e :=updateUser("USERNAME", Data.UserName, AccessToken)
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
