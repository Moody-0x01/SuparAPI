package main

import (
	"fmt"
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
const DefaultUserImg string = "/img/defUser.jpg"
const DefaultUserBg string = "/img/defBg.jpg"

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
    
    if res["code"] == 200 {
        return true, res["data"]["url"]
    }
    
    return false, res["data"]
    
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
    
    if res["code"] == 200 {
    	return true, res["data"]["url"]
    } 
    
    return false, res["data"]

}

func addPostImg_ToCDN(uuid string, Mime string, pid int) {
	
	values := make(map[string]interface{})
    
    values["id"] = uuid;
    values["mime"] = Mime;
    values["postID"] = pid;

    data, err := json.Marshal(values)
    
    resp, err := http.Post(addBG, "application/json" , bytes.NewBuffer(data))

    if err != nil {
        log.Fatal(err)
    };

    var res map[string]interface{};

    json.NewDecoder(resp.Body).Decode(&res)
    
    if res["code"] == 200 {
    	return true, res["data"]["url"]
    }

    
    return false, res["data"]
    
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