package networking
import (
    "net"
  	"os"
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