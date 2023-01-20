package networking
import (
    "net"
  	"os"
)




const (
	OK 				 	   = 200
	Created          	   = 201
	Accepted         	   = 202
	NoContent 		 	   = 204
	MovedPermanently 	   = 301
	MovedTemporarily 	   = 302
	NotModified      	   = 304
	BadRequest             = 400
	Unauthorized           = 401
	Forbidden              = 403
	NotFound               = 404
	InternalServerError    = 500
	NotImplemented         = 501
	BadGateway             = 502
	ServiceUnavailable     = 503
)

func GetCurrentMacAddress() string {
    
    host, err := os.Hostname()
    
    if err != nil {
		return ""
    } 

    addr, err := net.LookupIP(host)
    
    if err != nil {
		return ""
    } 
    
    return addr[2].String()
}


func GetCurrentMachineIp() string {
    
    host, err := os.Hostname()
    
    if err != nil {
		return ""
    } 

    addr, err := net.LookupIP(host)
    
    if err != nil {
	return ""
    } 
    
    return addr[1].String()
}