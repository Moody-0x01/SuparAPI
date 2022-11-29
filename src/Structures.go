package main


type User struct {
	Id_ 		 int `json:"id_"`
	Email 		 string `json:"Email"`
	UserName 	 string `json:"UserName"`
	PasswordHash string `json:"PasswordHash"`
	Token 		 string `json:"token"`
	Img 		 string `json:"img"`
	Bg 			 string `json:"bg"`
	Bio 		 string `json:"bio"`
	Address		 string `json:"addr"`
}


type ID struct {
	Id_ int `json:"id_"`
}

type Query struct {
	Query_ int `json:"query"`
}

type Post struct {
	Text string `json:"text"`
	Img	 string `json:"img"`
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








