package main;

import (
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
)

/*[ DONE ]*/

func GetAllPostsRoute(c *gin.Context) {
	All := GetAllPosts()
	
	var res Response = MakeServerResponse(200, All)

	c.JSON(http.StatusOK, res)
}

func getUserPostsRoute(c *gin.Context) {
	var id string = GetFieldFromContext(c, "id_")
	fmt.Println(id)
	var UserPosts []Post = getUserPostById(id);
	var res Response = MakeServerResponse(200, UserPosts)
	c.JSON(http.StatusOK, res)
}


func getUsersRoute(c *gin.Context) {
	var q string = GetFieldFromContext(c, "q")
	var Users []User;
	
	if q != "" {
		Users = getUsersByQuery(q)
	} else {
		Users = getUsers()
	}
	
	var res Response = MakeServerResponse(200, Users)

	c.JSON(http.StatusOK, res)
}

func getUserByIdRoute(c *gin.Context) {
	var uuid string = c.Param("uuid")
	var User User = getUserById(uuid)
	var res Response = MakeServerResponse(200, User)
	c.JSON(http.StatusOK, res)
}

func login(c *gin.Context) {
	var LoginForm UserLogin;
	c.BindJSON(&User);
	
	var resp Response

	if len(LoginForm.Token) > 0 {
		resp = AuthenticateUserJWT(LoginForm.Token)
	} else {
		if len(LoginForm.Password) > 0 && len(LoginForm.Email) > 0 {
			User, Ok := AuthenticateUserByEmailAndPwd(LoginForm.Password, LoginForm.Email)
			
			if Ok {
				resp = MakeServerResponse(200, User)
			} else {
				resp = MakeServerResponse(500, "Wrong password!")
			}

		} else {
			resp = MakeServerResponse(500, "Missing request attributes, Email or password not specified.")
		}
	}

	c.JSON(http.StatusOK, resp);
}


/*[ DONE ]*/

func SIGNUP(c *gin.Context) {
	/*
		
		Note NEXT TO IMPLEMENT:
		i create user object.
		ii Parse the request body to the user object.
		iii Add the user with AddUser(u User) Function
		iv if everything was okay and user was added successfully
			-> then Make the JWT token
			-> set user.Token send the data with 200 code.
			if not then:
			-> MakeServerResponse(500, "was not added because {REASON}")
	*/

	/* 
		BLUEPRINT:
			def signUp():	
				data = ConstructModel(request.json)
				Api_DB_HANDLER.connect()
				User = Api_DB_HANDLER.AddNewUser(data)
				if User.code == 200:
					AccessToken = User.data["T"]
					User = getUserByAT(AccessToken)
					User["Token"] = EncodeJWT(JWT_SECRET, {"T": AccessToken})
					return MakeServerResponse(200, User)
				return dumps(User.makeResponse())
			
	*/
	fmt.Println("Not Implemented")
}

