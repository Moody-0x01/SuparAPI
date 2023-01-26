package database;

import (
	"fmt"
	"github.com/Moody0101-X/Go_Api/models"
)
	
func Follow(follower_id int, followed_id int, Token string) models.Result {
	// "INSERT  INTO FOLLOWERS(follower_id, followed_id) VALUES(?, ?)"

	id, ok := GetUserIdByToken(Token)
	
	if ok {
		if id == follower_id {
			stmt, _ := DATABASE.Prepare("INSERT INTO FOLLOWERS(follower_id, followed_id) VALUES(?, ?)")
			_, err := stmt.Exec(follower_id, followed_id)

			if err != nil {
				fmt.Println("ERR: ", err)
				return models.MakeServerResult(false, "could not follow..")
			}
			
	    	Notification := models.NewNot(models.FOLLOW, followed_id, follower_id);
	    	pushNotificationForUser(Notification, " followed you!")
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

			stmt, _ := DATABASE.Prepare("DELETE FROM FOLLOWERS WHERE follower_id=? and followed_id=?")
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

func GetFollowers(followed int) []int {
	// "SELECT * FROM FOLLOWERS WHERE followed_id=? ORDER BY ID DESC"
	// people who is following followed.
	var followers []int;

	row, err := DATABASE.Query("SELECT follower_id FROM FOLLOWERS WHERE followed_id=? ORDER BY ID DESC", followed)
	
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

	row, err := DATABASE.Query("SELECT followed_id FROM FOLLOWERS WHERE follower_id=? ORDER BY ID DESC", following)
	
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

	row, err := DATABASE.Query("SELECT follower_id FROM FOLLOWERS WHERE follower_id=? AND followed_id=? ORDER BY ID DESC", follower, followed)
	
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

	return !(follower_id == 0);
}