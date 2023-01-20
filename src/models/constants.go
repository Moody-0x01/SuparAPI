package models;
import "github.com/Moody0101-X/Go_Api/networking"



var (
	FOLLOW                     = 0
    LIKE                       = 1
    COMMENT                    = 2
    NOTIFICATION               = "NN"
    NEWPOST                    = "NP"
    POSTCHANGE                 = "PC"
    MSG                        = "M"
    DefaultUserImg             = "/img/defUser.jpg"
    DefaultUserBg              = "/img/defBg.jpg"
    DefaultUserBio             = "Wait for it to load :)"
    DefaultUserAddress         = "Everywhere"
    CDN_API                    = "http://" + networking.GetCurrentMachineIp() + ":8500"
)