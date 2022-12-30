package database;

import (
	"fmt"
	"github.com/Moody0101-X/Go_Api/models"
	"github.com/Moody0101-X/Go_Api/cdn"
)

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


func Add_comment(uuid int, commentText string, PostId int, Token string, PostOwnerId int) models.Result {
	// ID INTEGER PRIMARY KEY AUTOINCREMENT,
    // uuid INTEGER,
    // post_id integer,
    // comment_text TEXT
    // TODO: add new id that specifies which user owns the post u liked or comminted on.

    id, ok := GetUserIdByToken(Token);

    if ok {

    	if id == uuid {
    		stmt, _ := dataBase.Prepare("INSERT INTO COMMENTS(uuid, post_id, comment_text) VALUES(?, ?, ?)")
			_, err := stmt.Exec(uuid, PostId, commentText)

			if err != nil {
				fmt.Println("ERR: ", err)
				return models.MakeServerResult(false, "could not add comment to db.")
			}
			
			Notification := models.NewNot(models.COMMENT, PostOwnerId, uuid);
	    	pushNotificationForUser(Notification, " commented on your post!")

			return models.MakeServerResult(true, "success")
    	}	
		
		return models.MakeServerResult(false, "token does not match this user, please make sure you are logged in.")

    }

    return models.MakeServerResult(false, "coult not get user id.")
	
}

func Add_like(uuid int, PostId int, Token string, PostOwnerId int) models.Result {
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
			
			var Notification models.Notification = models.NewNot(models.LIKE, PostOwnerId, uuid);
	    	pushNotificationForUser(Notification, " liked your post!")

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

















