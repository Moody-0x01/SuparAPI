package main;

import "fmt"

func db() {
	fmt.Println("Hi thhis is db!");
}

// TODO verifyUserPassword, AddUser, MakePost.

func print(s string) {
	fmt.Println(s);
}


/*-------------------------------------------------------------------------------------------------------------------------------
 	POSTS
-------------------------------------------------------------------------------------------------------------------------------*/


// func verifyUserPassword(LoginObject UserLogin) (User, bool) {
// 	// TODO	
// }


// func MakePost(post Post) bool {
// 	//TODO
// }

func GetAllPosts() []Post {
	var Posts []Post

	row, err := dataBase.Query("SELECT USER_ID, Text, IMG FROM POSTS ORDER BY ID DESC")
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return Posts
	}

	var temp Post

	for row.Next() {
		row.Scan(&temp.Uid_,&temp.Text, &temp.Img)
		temp.user = getUserById(temp.Uid_)
		Posts = append(Posts, temp)
	}

	// fmt.Println(Posts)
	return Posts
}

func getUserPostById(id string) []Post {
	// A functions to use 
	var Posts []Post

	row, err := dataBase.Query("SELECT Text, IMG FROM POSTS WHERE USER_ID=? ORDER BY ID DESC", id)
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return Posts
	}

	var temp Post

	for row.Next() {
		row.Scan(&temp.Text, &temp.Img);
		Posts = append(Posts, temp);
	}

	return Posts
}



func AuthenticateUserByEmailAndPwd(Pwd string, Email string) (User, Error) {
	
	var user User

	// row, err := dataBase.Query("SELECT ID, EMAIL, PASSWORDHASH, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE ID=? ORDER BY ID DESC", id)
	row, err := dataBase.Query("SELECT PASSWORDHASH FROM USERS WHERE EMAIL=?", Email)
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return user, MakeServerError(false, "Could not get user from db. 82")
	}

	var pwdHash string

	for row.Next() {
		row.Scan(&pwdHash)
	}
	
	if sha256_(Pwd) == pwdHash {
		//TODO Get user.
		//TODO Return user.
		row, err := dataBase.Query("SELECT ID, EMAIL, USERNAME, TOKEN, IMG, BG, BIO, ADDR FROM USERS WHERE EMAIL=? ORDER BY ID DESC", Email)

		defer row.Close()

		if err != nil {
			return user, MakeServerError(false, "Could not get user from db. 97")
		}

		for row.Next() {
			row.Scan(&user.Id_, &user.Email, &user.UserName, &user.Token, &user.Img, &user.Bg,  &user.Bio, &user.Address)
		}

		JWT, err := StoreTokenInJWT(user.Token)

		if err == nil {
			user.Token = JWT
			return user, MakeServerError(true, "User created! you can login now..")
		}

		var EmptyUser User
		return EmptyUser, MakeServerError(false, "Server had a problem encoding the token..")

	}
	
	return user, MakeServerError(false, "incorrect password. try again")
}

/*-------------------------------------------------------------------------------------------------------------------------------
 	USERS
-------------------------------------------------------------------------------------------------------------------------------*/
func getUserById(id string) User {
	
	var User User

	row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE ID=? ORDER BY ID DESC", id)
	defer row.Close()
	if err != nil {
		fmt.Println(err)
		return User
	}

	for row.Next() {
		row.Scan(&User.Id_, &User.UserName,&User.Img, &User.Bg, &User.Bio, &User.Address)
	}



	return User
}

func getUserByToken(Token string) (User, error) {
	var user User

	row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE TOKEN=? ORDER BY ID DESC", Token)
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return user, err
	}
	
	fmt.Println("Hi Getting user by ", Token)

	for row.Next() {
		row.Scan(&user.Id_, &user.UserName, &user.Img, &user.Bg, &user.Bio, &user.Address)
		// fmt.Println("Id_: ", user.Id_)
		// fmt.Println("UserName: ", user.UserName)
		// fmt.Println("Img: ", user.Img)
		// fmt.Println("Bg: ", user.Bg)
		// fmt.Println("Bio: ", user.Bio)
		// fmt.Println("Address: ", user.Address)
	}

	return user, nil
}

func getUsers() []User {
	var Users []User
	row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS ORDER BY ID DESC")
	defer row.Close()
	if err != nil {
		fmt.Println(err)
		return Users
	}

	var temp User

	for row.Next() {
		
		row.Scan(&temp.Id_, &temp.UserName, &temp.Img, &temp.Bg, &temp.Bio, &temp.Address)
		Users = append(Users, temp)
	}

	return Users	
}

func getUsersByQuery(Q string) []User {
	// TODO Verify this one works, sqlite to find alike attributes in row.
	// DONE.
	var Users []User
	var NewQ string = "%" + Q + "%"
	row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE USERNAME LIKE ? ORDER BY ID DESC", NewQ)

	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return Users
	}

	for row.Next() {
		var temp User
		row.Scan(&temp.Id_, &temp.UserName, &temp.Img, &temp.Bg, &temp.Bio, &temp.Address)
		Users = append(Users, temp)	
	}

	return Users
}

func GetUserByJWToken(JWToken string) (User, bool) {
	Token, isValid := GetTokenFromJwt(JWToken);

	if isValid {
		return GetUserByToken(Token), true
	} else {
		var EmptyUser User
		return EmptyUser, true
	}
}

func GetUserByToken(Token string) User {
	var User User;
	
	row, err := dataBase.Query("SELECT ID, EMAIL, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE TOKEN=? ORDER BY ID DESC", Token)
	defer row.Close()
	if err != nil {
		fmt.Println(err)
		return User
	}

	for row.Next() {
		row.Scan(&User.Id_, &User.Email, &User.UserName, &User.Img, &User.Bg, &User.Bio, &User.Address)
	}

	return User;
}


func CheckUser(Email string) bool {
		
	row, err := dataBase.Query("SELECT ID FROM USERS WHERE EMAIL=? ORDER BY ID DESC", Email)
	defer row.Close()
	var u []User;

	if err != nil {
		fmt.Println(err)
		return false
	}

	var temp User

	for row.Next() {	
		row.Scan(&temp.Id_, &temp.UserName, &temp.Img, &temp.Bg, &temp.Bio, &temp.Address)
		u = append(u, temp)	
	}

	if len(u) == 1 {
		return true
	} else if len(u) == 0 {
		return false
	} else {
		return true
	}
}

func AddUser(user User) Response {
	
	/* 
	- Check if the email associated with the request already exists.
		if it exists return an error code and a string.
		else just go on with process
	- try to add the user by:
		generating a token and assigning it to the user before adding
	*/

	

	if !CheckUser(user.Email) {

		var Token string = generateAccessToken(user.Email)
		stmt, _ := dataBase.Prepare("INSERT INTO USERS(EMAIL, USERNAME, PASSWORDHASH, TOKEN, IMG, BIO, BG, ADDR) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
		_, err := stmt.Exec(user.Email, user.UserName, user.PasswordHash, Token, user.Img, user.Bio, user.Bg, user.Address)
		
		if err != nil {
			return MakeServerResponse(500, "Could not add to db.")
		}

		FetchedUser, err := getUserByToken(Token)

		if err != nil {
			return MakeServerResponse(500, "Could not get created user from db. L288")
		} else {
			
			JWT, err := StoreTokenInJWT(Token)
			
			if err != nil {
				fmt.Println(err)
				return MakeServerResponse(500, "The server had a problem making the jwt token.")
			}

			fmt.Println("JWT: ", JWT)

			FetchedUser.Token = JWT;
			return MakeServerResponse(200, FetchedUser) // Success.
		}

	}

	return MakeServerResponse(500, "This user already exists..")
}


func updateUser(field string, newValue string, Token string) Error {
	
	var ok bool;
	var Query string

	switch field {
		case "IMG":
			Query = "UPDATE USERS SET IMG=? WHERE TOKEN=?"
			ok = true
			break		
		case "BIO":
			Query = "UPDATE USERS SET BIO=? WHERE TOKEN=?"
			ok = true
			break	
		case "ADDR":
			Query = "UPDATE USERS SET ADDR=? WHERE TOKEN=?"
			ok = true
			break		
		case "BG":
			Query = "UPDATE USERS SET BG=? WHERE TOKEN=?"
			ok = true
			break
		case "USERNAME":
			Query = "UPDATE USERS SET USERNAME=? WHERE TOKEN=?"
			ok = true
			break
		default:
			ok = false
			break
	}

	if ok {

		stmt, _ := dataBase.Prepare(Query)
		_, err := stmt.Exec(newValue, Token)

		if err != nil {
			return MakeServerError(false, "db err, could not update.")
		}	

		return MakeServerError(true, "success!")
	}

	return MakeServerError(false, "Unexpected field name.")
}


func AddPost(Text string, Img string, uuid int) Error {

	stmt, _ := dataBase.Prepare("INSERT INTO POSTS(USER_ID, Text, IMG) VALUES(?, ?, ?)")
	_, err := stmt.Exec(uuid, Text, Img)

	if err != nil {
		return MakeServerError(false, "could not add post. err L334")
	}

	return MakeServerError(true, "success!")
}

// TODO Add comment feeature.

// TODO Add likes feature.

// TODO Add FOLLOW/UNFOLLOW feature.