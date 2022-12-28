package cdn

import (
	// "fmt"
	"log"
    "io/ioutil"
    "net/http"
    // "net/url"
    "encoding/json"
    "bytes"
    "strings"
    "github.com/Moody0101-X/Go_Api/models"
)

var api string = GetCdnLink("./cdn.txt")
// const api string = "http://192.168.79.20:8500"
var addIMG string = api + "/Zimg/addAvatar"
var addBG string = api + "/Zimg/addbg"
var addPOST string = api + "/Zimg/NewPostImg"

func GetCdnLink(fname string) string {
    body, err := ioutil.ReadFile(fname)
    
    if err != nil {
        log.Fatalf("unable to read file: %v", err)
    }

    var next string = strings.TrimSpace(string(body))
    next = strings.Trim(next, "\n")
    return next
}

func AddUserAvatarToCdn(uuid int, Mime string) (bool, string) {

    if Mime == models.DefaultUserImg {
        return true, models.DefaultUserImg
    }

    values := make(map[string]interface{})
    
    values["id"] = uuid;
    values["mime"] = Mime;

    data, err := json.Marshal(values)

    resp, err := http.Post(addIMG, "application/json" , bytes.NewBuffer(data))

    if err != nil {
        log.Fatal(err)
    }

    var res map[string]interface{};

    json.NewDecoder(resp.Body).Decode(&res)
    
    if int(res["code"].(float64)) == 200 {
        return true, res["data"].(map[string]interface{})["url"].(string)
    }

    return false, res["data"].(string)
}

func AddUserBackgroundToCdn(uuid int , Mime string) (bool, string) {

    if Mime == models.DefaultUserBg {
        return true, models.DefaultUserBg
    }

    values := make(map[string]interface{})
    
    values["id"] = uuid;
    values["mime"] = Mime;

    data, err := json.Marshal(values)

    resp, err := http.Post(addBG, "application/json" , bytes.NewBuffer(data))

    if err != nil {
        log.Fatal(err)
    };

    var res map[string]interface{};

    json.NewDecoder(resp.Body).Decode(&res)
    
    if int(res["code"].(float64)) == 200 {
    	return true, res["data"].(map[string]interface{})["url"].(string)
    } 
    
    return false, res["data"].(string)
}

func AddPostImage(uuid int, Mime string, pid int) (bool, string) {
    
	values := make(map[string]interface{})
    
    values["id"] = uuid;
    values["mime"] = Mime;
    values["postID"] = pid;

    data, err := json.Marshal(values)
    
    resp, err := http.Post(addPOST, "application/json" , bytes.NewBuffer(data))

    if err != nil {
        log.Fatal(err)
    };

    var res map[string]interface{};

    json.NewDecoder(resp.Body).Decode(&res)
    
    if int(res["code"].(float64)) == 200 {
    	return true, res["data"].(map[string]interface{})["url"].(string)
    }
    
    return false, res["data"].(string)

}