package database;

import (
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/Moody0101-X/Go_Api/models"
	"github.com/Moody0101-X/Go_Api/crypto"
	"github.com/Moody0101-X/Go_Api/cdn"
)

var dataBase *sql.DB

// db initializer: Opens the db, then evluates a global conn variable.

func InitializeDb() (error, string) {
	
	var err error
	var dbPath string = "./db/Users.db"

	dataBase, err = sql.Open("sqlite3", dbPath); if err != nil {
		return err, ""
	}

	return nil, dbPath
}

func isEmpty(s string) bool { return len(s) == 0 }


func AuthenticateUserJWT(UserJWT string) models.Response {
   
    Token, Ok := crypto.GetTokenFromJwt(UserJWT)

    if Ok {
        fmt.Println("", crypto.Sha256_(Token))
        User_, err := GetUserByToken(Token)
        if err != nil {
            // a db error.
            return models.MakeServerResponse(500, "Db Error. (line 108).")
        } else {
            // Returns the user if everything was alright.
            return models.MakeServerResponse(200, User_)
        }

    } else {
        // JWT error
        return models.MakeServerResponse(500, "server could not decode the token. (line 117)")
    }
}
/*-------------------------------------------------------------------------------------------------------------------------------
 	POSTS
-------------------------------------------------------------------------------------------------------------------------------*/

func DeleteUserPost(PostId int, uuid int, Token string) models.Response {

	FetchedUser, err := GetUserByToken(Token);
	
	ownerId, ok := GetPostOwnerId(PostId);

	if ok {

		if uuid != ownerId { return models.MakeServerResponse(401, "Not authorized!") }
		if err != nil { return models.MakeServerResponse(500, "db error, could not fetch user by token.") }
		if uuid != FetchedUser.Id_ { return models.MakeServerResponse(401, "Not authorized!") }

		stmt, _ := dataBase.Prepare("DELETE FROM POSTS WHERE ID=?")

		result, err := stmt.Exec(PostId) // Execute query.


		if err != nil {
			return models.MakeServerResponse(500, "Could not delete the post.")
		}

		fmt.Println("", result);

		return models.MakeServerResponse(200, "success")
	}

	return models.MakeServerResponse(401, "Not authorized!")	
}
func GetAllPosts() []models.Post {
	var Posts []models.Post

	row, err := dataBase.Query("SELECT ID, USER_ID, Text, IMG FROM POSTS ORDER BY ID DESC")
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return Posts
	}

	var temp models.Post

	for row.Next() {
		row.Scan(&temp.Id_, &temp.Uid_,&temp.Text, &temp.Img)
		temp.User_ = GetUserById(temp.Uid_)
		Posts = append(Posts, temp)
	}
	
	// fmt.Println(Posts)
	return Posts
}
func GetPostOwnerId(PostID int) (int, bool) {
	var id int
	row, err := dataBase.Query("SELECT USER_ID FROM POSTS WHERE ID=? ORDER BY ID DESC", PostID)
	defer row.Close()
	
	if err != nil {
		fmt.Println(err)
		return 0, false
	}

	for row.Next() {
		row.Scan(&id);
	}
	return id, true
}
func GetUserPostById(id int) []models.Post {
	// A functions to use 
	var Posts []models.Post
	row, err := dataBase.Query("SELECT ID, Text, IMG FROM POSTS WHERE USER_ID=? ORDER BY ID DESC", id)
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return Posts
	}

	var temp models.Post

	for row.Next() {
		row.Scan(&temp.Id_, &temp.Text, &temp.Img);
		Posts = append(Posts, temp);
	}

	return Posts
}
func AuthenticateUserByEmailAndPwd(Pwd string, Email string) (models.User, models.Result) {
	var EmptyUser models.User

	if CheckUser(Email) {

		var user models.User
		row, err := dataBase.Query("SELECT PASSWORDHASH FROM USERS WHERE EMAIL=?", Email)
		
		defer row.Close()

		if err != nil {
			fmt.Println(err)
			return user, models.MakeServerResult(false, "Could not get user from db. 82")
		}

		var pwdHash string

		for row.Next() {
			row.Scan(&pwdHash)
		}
		
		if crypto.Sha256_(Pwd) == pwdHash {
			row, err := dataBase.Query("SELECT ID, EMAIL, USERNAME, TOKEN, IMG, BG, BIO, ADDR FROM USERS WHERE EMAIL=? ORDER BY ID DESC", Email)

			defer row.Close()

			if err != nil {
				return user, models.MakeServerResult(false, "Could not get user from db. 97")
			}

			for row.Next() {
				row.Scan(&user.Id_, &user.Email, &user.UserName, &user.Token, &user.Img, &user.Bg,  &user.Bio, &user.Address)
			}

			JWT, err := crypto.StoreTokenInJWT(user.Token)

			if err == nil {
				user.Token = JWT
				return user, models.MakeServerResult(true, "User created! you can login now..")
			}


			return EmptyUser, models.MakeServerResult(false, "Server had a problem encoding the token..")
		}
		
		return user, models.MakeServerResult(false, "incorrect password. try again")
	}

	return EmptyUser, models.MakeServerResult(false, "incorrect Email.. check and try again!")
}

/*-------------------------------------------------------------------------------------------------------------------------------
 	USERS
-------------------------------------------------------------------------------------------------------------------------------*/

func GetUserById(id int) models.AUser {
	
	var User models.AUser

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

func GetUserByToken(Token string) (models.User, error) {
	
	var User_ models.User
	row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE TOKEN=? ORDER BY ID DESC", Token)
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return User_, err
	}

	for row.Next() {
		row.Scan(&User_.Id_, &User_.UserName, &User_.Img, &User_.Bg, &User_.Bio, &User_.Address)
	}

	return User_, nil
}

func GetUsers(uuid interface{}) []models.AUser {
	var Users []models.AUser
	switch uuid.(type) {
		case int:
			row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS ORDER BY ID DESC")
			defer row.Close()
			
			if err != nil {
				fmt.Println(err)
				return Users
			}

			var temp models.AUser

			for row.Next() {
				row.Scan(&temp.Id_, &temp.UserName, &temp.Img, &temp.Bg, &temp.Bio, &temp.Address)
				temp.IsFollowed = IsFollowing(temp.Id_, uuid.(int))
				Users = append(Users, temp)
			}

			break

		case bool:

			row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS ORDER BY ID DESC")
			defer row.Close()
			
			if err != nil {
				fmt.Println(err)
				return Users
			}

			var temp models.AUser

			for row.Next() {
				
				row.Scan(&temp.Id_, &temp.UserName, &temp.Img, &temp.Bg, &temp.Bio, &temp.Address)
				Users = append(Users, temp)
			}

			break

		default:
			return Users
			break
	}
	

	return Users	
}

func GetUsersByQuery(Q string, uuid interface{}) []models.AUser {
	var Users []models.AUser
	var NewQ string = "%" + Q + "%"

	switch uuid.(type) {
		case int:
			row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE USERNAME LIKE ? ORDER BY ID DESC", NewQ)

			defer row.Close()

			if err != nil {
				fmt.Println(err)
				return Users
			}

			for row.Next() {
				var temp models.AUser
				row.Scan(&temp.Id_, &temp.UserName, &temp.Img, &temp.Bg, &temp.Bio, &temp.Address)
				temp.IsFollowed = IsFollowing(temp.Id_, int(uuid.(int)))
				Users = append(Users, temp)	
			}

			break

		case bool:
			
			row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE USERNAME LIKE ? ORDER BY ID DESC", NewQ)

			defer row.Close()

			if err != nil {
				fmt.Println(err)
				return Users
			}

			for row.Next() {
				var temp models.AUser
				row.Scan(&temp.Id_, &temp.UserName, &temp.Img, &temp.Bg, &temp.Bio, &temp.Address)
				Users = append(Users, temp)	
			}

			break

		default:
			return Users		
	}

	return Users
}

func GetUserByJWToken(JWToken string) (models.User, bool) {
	Token, isValid := crypto.GetTokenFromJwt(JWToken);

	if isValid {
		User, _ := GetUserByToken(Token)
		return User, true
	} else {
		var EmptyUser models.User
		return EmptyUser, true
	}
}

func GetUserIdByToken(t string) (int, bool) {
	var id int
	row, err := dataBase.Query("SELECT ID FROM USERS WHERE TOKEN=?", t)
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return id, false
	}

	for row.Next() {
		row.Scan(&id)
	}

	return id, true
}


func CheckUser(Email string) bool {
		
	row, err := dataBase.Query("SELECT ID FROM USERS WHERE EMAIL=? ORDER BY ID DESC", Email)
	defer row.Close()
	var u []models.User;

	if err != nil {
		fmt.Println(err)
		return false
	}

	var temp models.User

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

func GetNextUID(Table string) int {

	var id int;
	row, err := dataBase.Query("select MAX(ID) from " + Table);
	
	defer row.Close()

	if err != nil { return 0 }
	
	for row.Next() {
		row.Scan(&id);
	}

	return id + 1;
}

func AddUser(user models.User) models.Response {
	if !CheckUser(user.Email) {

		var Token string = crypto.GenerateAccessToken(user.Email)
		

		/*------------Add To cdn-------------*/
		var uuid = GetNextUID("Users")

		ok, img := cdn.AddUserAvatarToCdn(uuid, user.Img)

		if !ok {
			return models.MakeServerResponse(500, "cdn error, could not add avatar.")
		}

		ok, bg := cdn.AddUserBackgroundToCdn(uuid, user.Bg)

		if !ok {
			return models.MakeServerResponse(500, "cdn error, could not add background.")
		}


		/*------------Add To cdn-------------*/

		stmt, _ := dataBase.Prepare("INSERT INTO USERS(EMAIL, USERNAME, PASSWORDHASH, TOKEN, IMG, BIO, BG, ADDR) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
		_, err := stmt.Exec(user.Email, user.UserName, user.PasswordHash, Token, img, user.Bio, bg, user.Address)
		
		if err != nil {
			return models.MakeServerResponse(500, "Could not add to db.")
		}

		FetchedUser, err := GetUserByToken(Token)

		if err != nil {
			return models.MakeServerResponse(500, "Could not get created user from db. L288")
		} else {
			
			JWT, err := crypto.StoreTokenInJWT(Token)
			
			if err != nil {
				fmt.Println(err)
				return models.MakeServerResponse(500, "The server had a problem making the jwt token.")
			}

			FetchedUser.Token = JWT;
			return models.MakeServerResponse(200, FetchedUser) // Success.
		}

	}

	return models.MakeServerResponse(500, "This user already exists..")
}

func GetuuidByToken(Token string) (int, bool) {
	
	var uuid int
	row, err := dataBase.Query("SELECT ID FROM USERS WHERE TOKEN=?", Token)
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return 0, false
	}

	for row.Next() {
		row.Scan(&uuid);
	}

	return uuid, true
}

func UpdateUser(field string, newValue string, Token string) models.Result {
	
	var ok bool;
	var Query string
	

	switch field {
		case "IMG":
			Query = "UPDATE USERS SET IMG=? WHERE TOKEN=?"
			uuid, OK := GetuuidByToken(Token)
			
			if OK {
				OK, newValue = cdn.AddUserAvatarToCdn(uuid, newValue)
				fmt.Println("path: ", newValue)
				ok = OK
			}

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
			
			uuid, OK := GetuuidByToken(Token)
			
			if OK {
				OK, newValue = cdn.AddUserBackgroundToCdn(uuid, newValue)
				ok = OK
			}

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
			fmt.Println("db err: ", err)
			return models.MakeServerResult(false, "db err, could not update.")
		}

		return models.MakeServerResult(true, "success!")
	}

	return models.MakeServerResult(false, "Unexpected field name, or could not find user by token..")

}


func AddPost(Text string, Img string, uuid int) models.Result {

	if !isEmpty(Img) {
		var pid int = GetNextUID("Posts")
		
		var ok bool;
		
		ok, Img = cdn.AddPostImage(uuid, Img, pid) // (bool, string)
	
		if !ok {
			return models.MakeServerResult(false, "could not add post img to cdn.. err L480")
		}

	}

	stmt, _ := dataBase.Prepare("INSERT INTO POSTS(USER_ID, Text, IMG) VALUES(?, ?, ?)")
	_, err := stmt.Exec(uuid, Text, Img)

	if err != nil {
		return models.MakeServerResult(false, "could not add post. err L489")
	}

	return models.MakeServerResult(true, "success!")
}

// TODO Add comment feeature.

// TODO Add likes feature.
// TODO Add FOLLOW/UNFOLLOW feature.

func Add_comment(uuid int, commentText string, PostId int, Token string) models.Result {
	// ID INTEGER PRIMARY KEY AUTOINCREMENT,
    // uuid INTEGER,
    // post_id integer,
    // comment_text TEXT
    
    id, ok := GetUserIdByToken(Token)
    if ok {
    	if id == uuid {
    		stmt, _ := dataBase.Prepare("INSERT INTO COMMENTS(uuid, post_id, comment_text) VALUES(?, ?, ?)")
			_, err := stmt.Exec(uuid, PostId, commentText)

			if err != nil {
				fmt.Println("ERR: ", err)
				return models.MakeServerResult(false, "could not add comment to db.")
			}

			return models.MakeServerResult(true, "success")
    	}	
		
		return models.MakeServerResult(false, "token does not match this user, please make sure you are logged in.")

    }

    return models.MakeServerResult(false, "coult not get user id.")
	
}

func Add_like(uuid int, PostId int, Token string) models.Result {
	// ID INTEGER PRIMARY KEY AUTOINCREMENT,
	// uuid INTEGER,
	// post_id INTEGER

	id, ok := GetUserIdByToken(Token)

	if ok {
		if id == uuid {
			stmt, _ := dataBase.Prepare("INSERT INTO LIKES(uuid, post_id) VALUES(?, ?)")
			_, err := stmt.Exec(uuid, PostId)

			if err != nil {
				fmt.Println("ERR: ", err)
				return models.MakeServerResult(false, "could not add like to db.")
			}

			return models.MakeServerResult(true, "success")	
		}

		return models.MakeServerResult(false, "token does not match this user, please make sure you are logged in.")
	}

	return models.MakeServerResult(false, "coult not get user id.")
	
}

func Remove_like(uuid int, PostId int, Token string) models.Result {
	// ID INTEGER PRIMARY KEY AUTOINCREMENT,
	// uuid INTEGER,
	// post_id INTEGER

	id, ok := GetUserIdByToken(Token)

	if ok {
		if id == uuid {
			stmt, _ := dataBase.Prepare("DELETE FROM LIKES WHERE uuid=? AND post_id=?")
			_, err := stmt.Exec(uuid, PostId)

			if err != nil {
				fmt.Println("ERR: ", err)
				return models.MakeServerResult(false, "could not add like to db.")
			}

			return models.MakeServerResult(true, "success")
		}

		return models.MakeServerResult(false, "token does not match this user, please make sure you are logged in.")
	}

	return models.MakeServerResult(false, "coult not get user id.")
	
}


func Get_comments(PostId int) []models.Comment {
	// finished this one.
	// ID INTEGER PRIMARY KEY AUTOINCREMENT,
    // uuid INTEGER,
    // post_id integer,
    // comment_text TEXT
	
	var comments []models.Comment

	row, err := dataBase.Query("SELECT ID, uuid, comment_text FROM COMMENTS WHERE post_id=? ORDER BY ID DESC", PostId)
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return comments
	}

	var temp models.Comment

	for row.Next() {
		row.Scan(&temp.Id_, &temp.Uuid, &temp.Text);
		temp.User_ = GetUserById(temp.Uuid)
		comments = append(comments, temp)
	}

	return comments
}

func Get_likes(PostId int) []models.Like {

	/* 
	map[string]&websocket.conn
	map[uuid]&w
	*/

	var likes []models.Like

	row, err := dataBase.Query("SELECT ID, uuid FROM LIKES WHERE post_id=? ORDER BY ID DESC", PostId)

	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return likes
	}

	var temp models.Like

	for row.Next() {
	row.Scan(&temp.Id_,&temp.Uuid)
		temp.User_ = GetUserById(temp.Uuid)
		likes = append(likes, temp)
	}

	return likes
}

// GET Followers. by id.
// Add follower.
// 
// for fetch posts

// TODO Add follower, Notification+++
func Follow(follower_id int, followed_id int, Token string) models.Result {
	// "INSERT  INTO FOLLOWERS(follower_id, followed_id) VALUES(?, ?)"
	id, ok := GetUserIdByToken(Token)

	if ok {
		if id == follower_id {
			stmt, _ := dataBase.Prepare("INSERT INTO FOLLOWERS(follower_id, followed_id) VALUES(?, ?)")
			_, err := stmt.Exec(follower_id, followed_id)

			if err != nil {
				fmt.Println("ERR: ", err)
				return models.MakeServerResult(false, "could not follow..")
			}

			return models.MakeServerResult(true, "success")
		}

		return models.MakeServerResult(false, "token does not match this user, please make sure you are logged in.")
	}

	return models.MakeServerResult(false, "coult not get user id.")
}

func Unfollow(follower_id int, followed_id int, Token string) models.Result {
	// "DELETE FROM FOLLOWERS WHERE follower_id=? and followed_id=?"
	id, ok := GetUserIdByToken(Token)

	if ok {
		if id == follower_id {

			stmt, _ := dataBase.Prepare("DELETE FROM FOLLOWERS WHERE follower_id=? and followed_id=?")
			_, err := stmt.Exec(follower_id, followed_id)

			if err != nil {
				fmt.Println("ERR: ", err)
				return models.MakeServerResult(false, "could not unfollow.")
			}

			return models.MakeServerResult(true, "success")
		}

		return models.MakeServerResult(false, "token does not match this user, please make sure you are logged in.")
	}

	return models.MakeServerResult(false, "coult not get user id.")
}

// func pushNotification() {
// 	// Add Later.
// }

func GetFollowers(followed int) []int {
	// "SELECT * FROM FOLLOWERS WHERE followed_id=? ORDER BY ID DESC"
	// people who is following followed.
	var followers []int;

	row, err := dataBase.Query("SELECT follower_id FROM FOLLOWERS WHERE followed_id=? ORDER BY ID DESC", followed)
	
	//  CREATE TABLE FOLLOWERS (
 	//        ID INTEGER PRIMARY KEY AUTOINCREMENT,
 	//        followed_id INTEGER
 	//        follower_id INTEGER
	// );

	defer row.Close()

	if err != nil {
		fmt.Println("ERROR: \n", err)
		return followers
	}

	var uuid int
	var index = 0;
	
	for row.Next() {
		row.Scan(&uuid)
		followers = append(followers, uuid)
		index++
	}

	return followers
}


func GetFollowings(following int) []int {
	// people who followed is followingg..

	var followers []int;

	row, err := dataBase.Query("SELECT followed_id FROM FOLLOWERS WHERE follower_id=? ORDER BY ID DESC", following)
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return followers
	}

	var uuid int

	for row.Next() {
		row.Scan(&uuid)
		followers = append(followers, uuid)
	}

	return followers
}


func IsFollowing(followed int, follower int) bool {
	// "SELECT * FROM FOLLOWERS WHERE followed_id=? ORDER BY ID DESC"
	// people who followed is followingg..

	row, err := dataBase.Query("SELECT follower_id FROM FOLLOWERS WHERE follower_id=? AND followed_id=? ORDER BY ID DESC", follower, followed)
	
	//  CREATE TABLE FOLLOWERS (
 	//        ID INTEGER PRIMARY KEY AUTOINCREMENT,
 	//        followed_id INTEGER
 	//        follower_id INTEGER
	// );

	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return false
	}

	var follower_id int;

	for row.Next() {
		row.Scan(&follower_id)
		fmt.Println(follower_id)
	}

	return !(follower_id == 0)
}