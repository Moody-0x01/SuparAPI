package main

// Default fields for the user object.
var DefaultUserImg string = "/img/defUser.jpg"
var DefaultUserBg string = "/img/defBg.jpg"
var DefaultUserBio string = "Wait for it to load :)"
var DefaultUserAddress string = "Everywhere"


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


// for fetch posts
type Post struct {
	Id_  int 	`json:"id"`
	Uid_ int 	`json:"uuid"`
	Text string `json:"text"`
	Img	 string `json:"img"`
	user User   `json:"user"` 
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



type Error struct {
	Ok   bool `json:"ok"`
	Text string `json:"text"`
}

type ID struct {
	Id_ int `json:"id_"`
}

type Query struct {
	Query_ int `json:"query"`
}

// for securly adding posts.
type TokenizedPost struct {
	Token string `json:"token"`
	Uuid  int `json:"uuid"`
	Text  string `json:"text"`
	Img   string `json:"img"`
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
	
	var Resp Response;
	Resp.Code = code

	switch data.(type) {
		
		case []Post:
			Resp.Data = data.([]Post)
			break

		case []User:
			Resp.Data = data.([]User)
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
			Resp.Data = data.(string)
			break
	}

	return Resp
}

func MakeServerError(ok bool, t string) Error {
	var e Error
	e.Ok = ok
	e.Text = t
	return e
}






