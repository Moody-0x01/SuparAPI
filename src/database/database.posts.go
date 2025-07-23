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
		if uuid != ownerId { return models.MakeGenericServerResponse(401, "Not authorized!") }
		if err != nil { return models.MakeGenericServerResponse(500, "db error, could not fetch user by token.") }
		if uuid != FetchedUser.Id_ { return models.MakeGenericServerResponse(401, "Not authorized!") }

		stmt, _ := DATABASE.Prepare("DELETE FROM POSTS WHERE ID=?")
		_, err := stmt.Exec(PostId) // Execute query.


		if err != nil {
			return models.MakeGenericServerResponse(500, "Could not delete the post.")
		}

		return models.MakeGenericServerResponse(200, "success")
	}

	return models.MakeGenericServerResponse(401, "Not authorized!")	
}

func GetAllPosts() map[int]models.Post {
	
	Posts := make(map[int]models.Post);
	row, err := DATABASE.Query("SELECT ID, USER_ID, Text, IMG, CREATED_DATE FROM POSTS ORDER BY ID DESC")
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return Posts
	}

	var tempPost models.Post

	for row.Next() {
		row.Scan(&tempPost.Id_, &tempPost.Uid_,&tempPost.Text, &tempPost.Img, &tempPost.Date);
		tempPost.User_ = GetUserById(tempPost.Uid_);
		tempPost.Img = CheckCdnLink(tempPost.Img); // Fix for the api cdn images...
		
		tempPost.PostLikes = Get_likes(tempPost.Id_);
		tempPost.LikesCount = len(tempPost.PostLikes);
		tempPost.PostComments = Get_comments(tempPost.Id_)
		tempPost.CommentsCount = len(tempPost.PostComments);
		
		Posts[tempPost.Id_] = tempPost;
	}

	return Posts;
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

func GetUserPostById(id int) map[int]models.Post {
	// A functions to use 
	Posts := make(map[int]models.Post);
	
	row, err := DATABASE.Query("SELECT ID, Text, IMG, CREATED_DATE FROM POSTS WHERE USER_ID=? ORDER BY ID DESC", id)
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return Posts
	}

	var tempPost models.Post

	for row.Next() {
		row.Scan(&tempPost.Id_, &tempPost.Text, &tempPost.Img, &tempPost.Date);
		tempPost.Img = CheckCdnLink(tempPost.Img);

		tempPost.PostLikes = Get_likes(tempPost.Id_);
		tempPost.LikesCount = len(tempPost.PostLikes);
		tempPost.PostComments = Get_comments(tempPost.Id_)
		tempPost.CommentsCount = len(tempPost.PostComments);
		
		Posts[tempPost.Id_] = tempPost;
	}

	return Posts
}

func GetPostById(Post_id int) models.Post {
	
	row, err := DATABASE.Query("SELECT ID, Text, IMG, USER_ID, CREATED_DATE FROM POSTS WHERE ID=? ORDER BY ID DESC", Post_id)
	
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
		PostOB.PostLikes = Get_likes(Post_id);
		PostOB.LikesCount = len(PostOB.PostLikes);
		PostOB.PostComments = Get_comments(Post_id)
		PostOB.CommentsCount = len(PostOB.PostComments);
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

	stmt, _ := DATABASE.Prepare("INSERT INTO POSTS(USER_ID, TEXT, IMG, CREATED_DATE) VALUES(?, ?, ?, datetime())")
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
	PostObj.PostLikes = Get_likes(PostObj.Id_);
	PostObj.LikesCount = len(PostObj.PostLikes);
	PostObj.PostComments = Get_comments(PostObj.Id_)
	PostObj.CommentsCount = len(PostObj.PostComments);
	PostObj.User_ = GetUserById(PostObj.Uid_);

	SockMsg := PostObj.EncodeToSocketResponse();
	models.ClientPool.BroadCastJSON(SockMsg, uuid)
	return models.MakeServerResult(true, pid)
}


// TODO func Remove_comment(commentId).
// TODO Send the comment/Like event to Users.
func Add_comment(uuid int, commentText string, PostId int, Token string, PostOwnerId int) models.Result {
/*
type Comment struct {
	Id_          int `json:"id_"`
	Post_id		 int `json:"post_id"`
	Uuid		 int `json:"uuid"`
	Text		 string `json:"text"`
	User_		 AUser `json:"user"` // Filled when fitching comments.
}

type Like struct {
	Id_          int `json:"id_"`
	Post_id		 int `json:"post_id"`
	Uuid		 int `json:"uuid"`
	User_		 AUser `json:"user"` // Filled when fitching comments.
}

*/
    id, ok := GetUserIdByToken(Token);

    if ok {

    	if id == uuid {
    		stmt, _ := DATABASE.Prepare("INSERT INTO COMMENTS(ID, POST_ID, COMMENT_TEXT) VALUES(?, ?, ?)")
			_, err := stmt.Exec(uuid, PostId, commentText)

			if err != nil {
				fmt.Println("ERR: ", err)
				return models.MakeServerResult(false, "could not add comment to db.")
			}
			
			Notification := models.NewNot(models.COMMENT, PostOwnerId, uuid);
			Notification.Post_id = PostId;
	    	pushNotificationForUser(Notification, " commented on your post!")

	    	var c models.Comment;
			
			c.Id_ = GetLastID("COMMENTS");
			c.Post_id = PostId;
			c.Uuid = uuid;
			c.Text = commentText;
			c.User_ = GetUserById(uuid);
			
			Message := c.EncodeToSocketResponse();
			models.ClientPool.BroadCastJSON(Message, uuid);

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
			stmt, _ := DATABASE.Prepare("INSERT INTO LIKES(USER_ID, POST_ID) VALUES(?, ?)")
			_, err := stmt.Exec(uuid, PostId)

			if err != nil {
				fmt.Println("ERR: ", err)
				return models.MakeServerResult(false, "could not add like to db.")
			}
			
			var Notification models.Notification = models.NewNot(models.LIKE, PostOwnerId, uuid);
			Notification.Post_id = PostId;
	    	pushNotificationForUser(Notification, " liked your post!")
	    	
	    	var like models.Like;

			like.Id_ = GetLastID("LIKES");
			like.Post_id = PostId;
			like.Uuid = uuid;
			like.User_ = GetUserById(uuid);
			
			Message := like.EncodeToSocketResponse();
			models.ClientPool.BroadCastJSON(Message, uuid);

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
			stmt, _ := DATABASE.Prepare("DELETE FROM LIKES WHERE USER_ID=? AND POST_ID=?")
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

	row, err := DATABASE.Query("SELECT ID, USER_ID, COMMENT_TEXT FROM COMMENTS WHERE POST_ID=? ORDER BY ID DESC", PostId)
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return comments
	}

	var tempComment models.Comment

	for row.Next() {
		row.Scan(&tempComment.Id_, &tempComment.Uuid, &tempComment.Text);
		tempComment.User_ = GetUserById(tempComment.Uuid)
		comments = append(comments, tempComment)
	}

	return comments
}

func Get_likes(PostId int) []models.Like {

	/* 
		map[string]&websocket.conn
		map[uuid]&w
	*/

	var likes []models.Like

	row, err := DATABASE.Query("SELECT ID, USER_ID FROM LIKES WHERE POST_ID=? ORDER BY ID DESC", PostId)

	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return likes
	}

	var tempLike models.Like
	for row.Next() {

	row.Scan(&tempLike.Id_,&tempLike.Uuid)
		tempLike.User_ = GetUserById(tempLike.Uuid)
		likes = append(likes, tempLike)
	}

	return likes
}
