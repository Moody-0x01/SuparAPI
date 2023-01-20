package models
import "time"

type Post struct {
	Id_   int 	`json:"id"`
	Uid_  int 	`json:"uuid"`
	Text  string `json:"text"`
	Img	  string `json:"img"`
	User_ AUser   `json:"user"` 
	Date  time.Time `json:"date"`
}

type TokenizedPost struct {
	PostID int `json:"id_"`
	Token string `json:"token"`
	Uuid  int `json:"uuid"`	
	Text  string `json:"text"`
	Img   string `json:"img"`
}

func (P *Post) EncodeToSocketResponse() SocketMessage { 
	return MakeSocketResp(NEWPOST, 200, P) 
}
