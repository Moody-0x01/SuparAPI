package models

import (
	"fmt"
)

func isEmpty(s string) bool { return len(s) == 0 }

type Result struct {
	Ok   bool `json:"ok"`
	Data interface{} `json:"data"`
}

type Query struct {
	Query_ int `json:"query"`
}

type Response struct {
	Code int `json:"code"`
 	Data interface{} `json:"data"`
}


type SocketMessage struct {
	Code   int         `json:"code"`
	Action int         `json:"action"`
	Data   interface{} `json:"data"`
}


func MakeGenericServerResponse(code int, data interface{}) Response {
	return Response{
		Code: code,
		Data: data,
	}
}

func MakeSocketResp(Enum int, code int, data interface{}) SocketMessage {
	
	var resp SocketMessage;
	resp.Code = code;
	resp.Action = Enum;
	
	switch data.(type) {
		
		case Notification:
			resp.Data = data.(Notification)
			break
		
		case Like:
			resp.Data = data.(Like)
			break
		
		case Comment:
			resp.Data = data.(Comment)
			break
		
		case Post:
			resp.Data = data.(Post)
			break

		case UMessage:
			resp.Data = data.(UMessage)
			break

		default:
			resp.Data = data
			break
	}

	return resp;
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
		case []Discussion:
			Resp.Data = data.([]Discussion)	
			break

		case Discussion:
			Resp.Data = data.(Discussion)
			break

		case []UMessage:
			Resp.Data = data.([]UMessage)
			break

		case UMessage:
			Resp.Data = data.(UMessage)
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

		case string:
			Resp.Data = data.(string)
			break
			
		default:
			fmt.Println("Unexpected data type. make sure it is added in MakeServerResponse(code int, data interface{}){ }")
			break
	}

	return Resp
}

func MakeServerResult(ok bool, t interface{}) Result {
	var e Result
	e.Ok = ok
	
	switch t.(type) {
		case int:
			e.Data = t.(int)
			break
		default:
			e.Data = t.(string)
			break
	}

	return e
}
