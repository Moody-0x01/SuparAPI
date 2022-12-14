package main

import (
	// "fmt"
	"log"
    // "io/ioutil"
    "net/http"
    // "net/url"
    "encoding/json"
    "bytes"
)


const api string = "http://localhost:8500"
const addIMG string = api + "/Zimg/addAvatar"
const addBG string = api + "/Zimg/addbg"
const addPOST string = api + "/Zimg/NewPostImg"

func addAvatar_ToCDN(uuid int, Mime string) (bool, string) {

    if Mime == DefaultUserImg {
        return true, DefaultUserImg
    }

    values := make(map[string]interface{})
    
    values["id"] = uuid;
    values["mime"] = Mime;

    data, err := json.Marshal(values)

    resp, err := http.Post(addIMG, "application/json" , bytes.NewBuffer(data))

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

func addbackground_ToCDN(uuid int , Mime string) (bool, string) {
    
    if Mime == DefaultUserBg {
        return true, DefaultUserBg
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

func addPostImg_ToCDN(uuid int, Mime string, pid int) (bool, string) {
	
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

/*
PYTHON VERSION:

	# CDN link
	api = "http://localhost:8500"

	# Endpoints.
	addIMG = f"{api}/Zimg/addAvatar"
	addBG = f"{api}/Zimg/addbg"
	addPOST = f"{api}/Zimg/NewPostImg"


	def addAvatar(uuid: int | str, Mime: str) -> dict:
	    
	    res = post(addIMG, json={
	        "id": uuid,
	        "mime": Mime
	    })

	    return res.json()

	def addbg(uuid: int | str, Mime: str) -> dict:
	    res = post(addBG, json={
	        "id": uuid,
	        "mime": Mime
	    })

	    return res.json()

	def addPost(uuid, Mime, postid=1):
	    res = post(addPOST, json={
	        "id": uuid,
	        "mime": Mime,
	        "postID": postid
	    })

	    return res.json()


*/