package database;

import (
	"fmt"
	"time"
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

		stmt, _ := DATABASE.Prepare("DELETE FROM POSTS WHERE ID=?")
		_, err := stmt.Exec(PostId) // Execute query.


		if err != nil {
			return models.MakeServerResponse(500, "Could not delete the post.")
		}

		return models.MakeServerResponse(200, "success")
	}

	return models.MakeServerResponse(401, "Not authorized!")	
}

func GetAllPosts() []models.Post {
	var Posts []models.Post

	row, err := DATABASE.Query("SELECT ID, USER_ID, Text, IMG, CreatedDate FROM POSTS ORDER BY ID DESC")
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return Posts
	}

	var temp models.Post

	for row.Next() {
		row.Scan(&temp.Id_, &temp.Uid_,&temp.Text, &temp.Img, &temp.Date)
		temp.User_ = GetUserById(temp.Uid_)
		temp.Img = CheckCdnLink(temp.Img); // Fix for the api cdn images...
		Posts = append(Posts, temp)
	}
	
	// fmt.Println(Posts)
	return Posts
}

func GetPostOwnerId(PostID int) (int, bool) {
	var id int
	row, err := DATABASE.Query("SELECT USER_ID FROM POSTS WHERE ID=? ORDER BY ID DESC", PostID)
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
	row, err := DATABASE.Query("SELECT ID, Text, IMG, CreatedDate FROM POSTS WHERE USER_ID=? ORDER BY ID DESC", id)
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return Posts
	}

	var temp models.Post

	for row.Next() {
		row.Scan(&temp.Id_, &temp.Text, &temp.Img, &temp.Date);
		temp.Img = CheckCdnLink(temp.Img);
		Posts = append(Posts, temp);
	}

	return Posts
}

func GetPostById(Post_id int) models.Post {
	
	row, err := DATABASE.Query("SELECT ID, Text, IMG, USER_ID, CreatedDate FROM POSTS WHERE ID=? ORDER BY ID DESC", Post_id)
	
	defer row.Close()
	
	var PostOB models.Post;
	
	if err != nil {
		fmt.Println(err)
		return PostOB
	}

	for row.Next() {
		row.Scan(&PostOB.Id_, &PostOB.Text, &PostOB.Img, &PostOB.Uid_, &PostOB.Date);
		PostOB.Img = CheckCdnLink(PostOB.Img);
		PostOB.User_ = GetUserById(PostOB.Uid_);
	}

	return PostOB
}

func AddPost(Text string, Img string, uuid int) (models.Result) {
	/* 

		TODO: 
			We have a bug here!, the uuid that is sent to the user afte the post is added is incorrect!
			so it messes up the rendering of the react comp tree..

			Because When we want to delete a certain post we need to verify if it is Your post!
			also, the id is the key of the UI post comp, so.. making alot of problems in the tree?

			Solution Idea1 has failed -> GetNextUID
			Solved, what a dumbass I was getting the maximum id before actually adding a new Post!!!!!!!!!!

	*/

	var pid int = -1;

	if !isEmpty(Img) {
		var ok bool;
		
		ok, Img = cdn.AddPostImage(uuid, Img, pid) // (bool, string)
	
		if !ok {
			return models.MakeServerResult(false, "could not add post img to cdn.. err L480")
		}

	}

	stmt, _ := DATABASE.Prepare("INSERT INTO POSTS(USER_ID, Text, IMG, CreatedDate) VALUES(?, ?, ?, datetime())")
	_, err := stmt.Exec(uuid, Text, Img)

	if err != nil { return models.MakeServerResult(false, "could not add post. err L149") }

	pid = GetNewPostID();
	
	if pid == 0 { return models.MakeServerResult(false, "could not get pid. err L154") }

	// TODO: Broadcast the post Msg...
	
	var PostObj models.Post;
	PostObj.Id_ = pid;
	PostObj.Uid_ = uuid;
	PostObj.Text = Text;
	PostObj.Img = CheckCdnLink(Img); 
	PostObj.Date = time.Now();

	PostObj.User_ = GetUserById(PostObj.Uid_);

	SockMsg := PostObj.EncodeToSocketResponse();
	
	models.ClientPool.BroadCastJSON(SockMsg, uuid)

	return models.MakeServerResult(true, pid)
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
    		stmt, _ := DATABASE.Prepare("INSERT INTO COMMENTS(uuid, post_id, comment_text) VALUES(?, ?, ?)")
			_, err := stmt.Exec(uuid, PostId, commentText)

			if err != nil {
				fmt.Println("ERR: ", err)
				return models.MakeServerResult(false, "could not add comment to db.")
			}
			
			Notification := models.NewNot(models.COMMENT, PostOwnerId, uuid);
			Notification.Post_id = PostId;
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
			stmt, _ := DATABASE.Prepare("INSERT INTO LIKES(uuid, post_id) VALUES(?, ?)")
			_, err := stmt.Exec(uuid, PostId)

			if err != nil {
				fmt.Println("ERR: ", err)
				return models.MakeServerResult(false, "could not add like to db.")
			}
			
			var Notification models.Notification = models.NewNot(models.LIKE, PostOwnerId, uuid);
			Notification.Post_id = PostId;
	    	pushNotificationForUser(Notification, " liked your post!")
	    	fmt.Println("Notification was pushed..")
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
			stmt, _ := DATABASE.Prepare("DELETE FROM LIKES WHERE uuid=? AND post_id=?")
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

	row, err := DATABASE.Query("SELECT ID, uuid, comment_text FROM COMMENTS WHERE post_id=? ORDER BY ID DESC", PostId)
	
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

	row, err := DATABASE.Query("SELECT ID, uuid FROM LIKES WHERE post_id=? ORDER BY ID DESC", PostId)

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