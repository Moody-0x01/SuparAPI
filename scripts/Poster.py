"""
This python script is a unit testing entity, I use it to test routes and make sure that everything works as expected.
ROUTE:
    
    router.POST("/login", login) expects => {
        "Email": v, 
        "Password": v
    } or {
        "T": v
    }

    router.POST("update", update) expects => {
        "token": v(important.), 
        "img": v | null, 
        "bg": v | null, 
        "bio": v | null, 
        "addr": v | null
    }

	router.POST("/NewPost", NewPost) expects => {
		"Token": v,
		"uuid": v,
		"text": v,
		"img": v
	}

	router.POST("/signup", signUp) expects => {
        "Email": v,
        "Password": v,
        "UserName": v
    }
	
    // Get routes.
    router.GET("/getUserPosts", getUserPostsRoute) expects => "/getUserPosts?id_={:UUID}"
    router.GET("/GetAllPosts", GetAllPostsRoute) expects none, returns all posts.
    router.GET("/query", getUsersRoute) expects "/query?q{search_query_as string.}"
    router.GET("/:uuid", getUserByIdRoute) expects /UserID returns user object.
    
"""

from requests import post
from base64 import b64encode
GUKU_TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJUIjoiNzM3Y2UxZWMwYmZkZmQyNDUyMWVhZWUxNzU0YzE4YWQ2MWRmZWYxYTNjMGJhOTcyNjBmZTE5ZWY4MGQwYTJlNiJ9.A6BKhdwW52KXRRO-RS6nKR46o0l5c6KNLXsq2xkpbD0'
UUID = 1714
Url = "http://localhost:8888/"
sign_up = f"{Url}signup"
login = f"{Url}login"
addPost = f"{Url}NewPost"
update = f"{Url}update"

TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJUIjoiTkdGbE1EVmhZamRrWkRZM09UUmhaV1l5TVRGaU5EaG1aakEzT0RJMFpHSmtOVFJpTURsak1EY3pNalExTVRoa1lUVXlNamsxWkdJMk1ESTBaalZrWlE9PSJ9.06MJ_eLYYPQmw1ekgn2pBG0JaO-zgc8IYE2vbo7S-Mw"

def MakeMime(fp: any): 
	ext = fp.name.split(".")[1]
	enc = b64encode(fp.read()).decode()
	return f"data:image/{ext};base64,{enc}"

def loadImgTob64():
    print("]![ Loading image")
    with open("../imagesTest/bug.jpg", "rb") as fp:
        return MakeMime(fp)

def updateAVATAR():
    """ TESTED: Success code returned. """
    
    return post(update, json={
        "img": loadImgTob64(), 
        "token": TOKEN
    }).json()

def addUserPost():
    """ TESTED: Success code returned. """
    data = {
		"Token": TOKEN,
		"uuid": UUID,
		"text": "Hello, this is made with python.",
        "img": ""
	}

    return post(addPost, json=data).json()

def addUser(data: dict) -> dict:
    """ TESTED: Success jwt returned. """
    return post(sign_up, json=data).json()

def loginU():
    """ TESTED: Success user returned. """
    return post(login, json={"T": GUKU_TOKEN}).json()

def  main():
    res = loginU()
    print(res)

if __name__ == "__main__": main()
