package models

import (
	"fmt"
)

func isEmpty(s string) bool { return len(s) == 0 }

type Post struct {
	Id_  int 	`json:"id"`
	Uid_ int 	`json:"uuid"`
	Text string `json:"text"`
	Img	 string `json:"img"`
	User_ AUser   `json:"user"` 
}



type Result struct {
	Ok   bool `json:"ok"`
	Text string `json:"text"`
}

type Query struct {
	Query_ int `json:"query"`
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

		case []Notification:
			Resp.Data = data.([]Notification)
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
