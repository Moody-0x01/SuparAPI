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

/*[ DONE ]*/

/* AUTH. */

func signup(c *gin.Context) {
}

func login(c *gin.Context) {
	
	var User UserLogin;
	
	c.BindJSON(&User);
	
	// TODO Decode requested data to a User Struct.
	/*
	    TODO Check if there is Token, 
	    if there is then:
		TODO Verify the token.
		TODO extract the token data we need to find the current user's data.
	    else:
		TODO Authenticate the User and return the user.
		TODO return response.. {"data": ..., "code": 200|500|201|404}
		TODO 
	*/

	fmt.Println(User.Email);
	fmt.Println(User.Password);
	var resp Response;
	resp.Code = 200;
	resp.Data = "Hii It was received.";
	c.JSON(http.StatusOK, resp);


	// if "AccessToken" in request.json:
	// 		Jwt_token = request.json["AccessToken"]
	// 		data = VerifyJWT(JWT_SECRET, Jwt_token)
	// 		if isinstance(data, dict):
	// 			# Try to get user data to send to my client.
	// 			accessToken = data['T']
	// 			User = getUserByAT(accessToken)
	// 			return MakeServerResponse(200, User)
	// 		elif isinstance(data, str):
	// 			# Return an  error code and the error message to the client.
	// 			return MakeServerResponse(202, data)

	// 	data = ConstructModel(request.json)
	// 	DB_HANDLER = createNewDbObject()
	// 	DB_HANDLER.connect()
	// 	User = DB_HANDLER.AuthenticateUser(data)
	// 	if User.code == 200:
	// 		response = User.makeResponse()
			
	// 		AccessToken = {
	// 			"T": User.data["Token"]
	// 		}

	// 		User.data["Token"] = EncodeJWT(JWT_SECRET, AccessToken)
	// 		return dumps(response)

	// 	return dumps(User.makeResponse())
}


/* Data base interactors. */











