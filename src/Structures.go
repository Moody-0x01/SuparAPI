package main
import "fmt"
// Default fields for the user object.
const DefaultUserImg string = "/img/defUser.jpg"
const DefaultUserBg string = "/img/defBg.jpg"
const DefaultUserBio string = "Wait for it to load :)"
const DefaultUserAddress string = "Everywhere"


type User struct {
	Id_ 		 int `json:"id_"`
	Email 		 string `json:"Email"`
	UserName 	 string `json:"UserName"`
	PasswordHash string `json:"Password"`
	Token 		 string `json:"token"`
	Img 		 string `json:"img"`
	Bg 			 string `json:"bg"`
	Bio 		 string `json:"bio"`
	Address		 string `json:"addr"`
}


// type Notification struct {
//     TYPE TEXT DEFAULT null, [follow | like | comment | ...]
//     USER_ID INTEGER,
//     OTHER_ID INTEGER,
//     PID INTEGER,
//     MSG TEXT DEFAULT null
//     ...

// }

type AUser struct {
	Id_ 		 int `json:"id_"`
	UserName 	 string `json:"UserName"`
	Img		 	 string `json:"img"`
	Bg 			 string `json:"bg"`
	Bio 		 string `json:"bio"`
	Address		 string `json:"addr"`
	IsFollowed	 bool `json:"isfollowed"`
}

type WSocketAccessController struct {
	Uuid 	int `json:"uuid"`
}

// type userConn struct {
// 	Connection string
// 	uuid string
// 	Connectionid_ string
// }
// type clients struct {
// 	clients []userConn
// }
// func (*clients c) sendMsg() {
// 	b := make([]byte, 1024)
	
// 	for i := 0; i < len(clients); i++ {
	
// 		n, err := clients[i].sendbytes(b)
// 		if err != nil {
// 			return 
// 		}
// 	}	
// }

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

type TFollow struct {
	Id_        		int `json:"id_"`
	Follower_id		int `json:"follower_id"`
	Followed_id		int `json:"followed_id"`
	UToken			string `json:"token"`
}

type Follow struct {
	Id_        		int `json:"id_"`
	Follower_id		int `json:"follower_id"`
	Followed_id		int `json:"follower_id"`
}

type Post struct {
	Id_  int 	`json:"id"`
	Uid_ int 	`json:"uuid"`
	Text string `json:"text"`
	Img	 string `json:"img"`
	User_ AUser   `json:"user"` 
}


func (U *User) setDefaults() {
	//TODO Setting the default fields to add to the db if some are not present.
	
	/*
		THOSE ARE THE FIELDS TO BE CHANGED if they were not set.
			Img 		 string `json:"img"`
			Bg 			 string `json:"bg"`
			Bio 		 string `json:"bio"`
			Address		 string `json:"addr"`

	*/

	if isEmpty(U.Img) {
		U.Img = DefaultUserImg
	}
	
	if isEmpty(U.Bg) {
		U.Bg = DefaultUserBg
	}
	
	if isEmpty(U.Bio) {
		U.Bio = DefaultUserBio
	}

	if isEmpty(U.Address) {
		U.Address = DefaultUserAddress
	}
}



type Result struct {
	Ok   bool `json:"ok"`
	Text string `json:"text"`
}

type ID struct {
	Id_ int `json:"id_"`
}

type Query struct {
	Query_ int `json:"query"`
}

// for securly adding or edit a posts.
type TokenizedPost struct {
	PostID int `json:"id_"`
	Token string `json:"token"`
	Uuid  int `json:"uuid"`
	
	Text  string `json:"text"`
	Img   string `json:"img"`
}

type TokenizedComment struct {
	Post_id		 int `json:"post_id"`
	Uuid		 int `json:"uuid"`
	Text		 string `json:"text"`
	Token        string `json:"token"`
}

type TokenizedLike struct {
	Post_id		 int `json:"post_id"`
	Uuid		 int `json:"uuid"`
	Token        string `json:"token"`
}

type UserLogin struct {
	Password string `json:"Password"`
	Email    string `json:"Email"`
	Token 	 string `json:"T"`
}


type Response struct {
	Code int `json:"code"`
 	Data interface{} `json:"data"`
}

func MakeServerResponse(code int, data interface{}) Response {
	var Resp Response
	Resp.Code = code

	switch data.(type) {
		
		case []Post:
			Resp.Data = data.([]Post)
			break
		case []Like:
			Resp.Data = data.([]Like)
			break

		case []Comment:
			Resp.Data = data.([]Comment)
			break

		case []User:
			Resp.Data = data.([]User)
			break

		case []AUser:
			Resp.Data = data.([]AUser)
			break
		case []int:
			Resp.Data = data.([]int)
			break

		case int:
			Resp.Data = data.(int)
			break

		case Like:
			Resp.Data = data.(Like)
			break

		case Comment:
			Resp.Data = data.(Comment)
			break

		case AUser:
			Resp.Data = data.(AUser)
			break

		case User:
			Resp.Data = data.(User)
			break

		case Post:
			Resp.Data = data.(Post)
			break

		case UserLogin:
			Resp.Data = data.(UserLogin)
			break


		default:
			fmt.Println("Unexpected data type. make sure it is added in MakeServerResponse(code int, data interface{}){ }")
			break
	}

	return Resp
}

func MakeServerResult(ok bool, t string) Result {
	var e Result
	e.Ok = ok
	e.Text = t
	return e
}
