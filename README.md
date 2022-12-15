

# Go Api
### Routes
NOTE: This section will be updated regularily once I add or remove a route.

- `/login` Login, can get a token or (password, Email) and if everything is valid it returns user info + Token.
- `/getUserPosts` return the user's posts, and it gets an id_ form value as uuid.
- `/GetAllPosts` return all the posts from db.
- `/query` used to query a specific user. if not found then an empty list is returned or an error response code.
- `/:uuid` quick user lookup by uuid.
- `/signup` making an account. then returning the token to update, delete, add and other operation regarding ur account.
- `/update` expects either a img, bio, addr or bg to be updated in the database, but also accepts a token that will be given if authenticated. if the token is not provided, the server will return error code 500.
- `/NewPost` add new post with token, expects a token, img, uuid, and post text. if something is not set properly it will return a response report about the error.
- `/DeletePost` deleting user post using userid, postid and jwt.

### Database (Tables and schema.)

I have used only 2 tables just because I did not want to overcomplicate things, but more will be added as m progressing in this project.
```sql

    CREATE TABLE USERS (
        ID INTEGER PRIMARY KEY AUTOINCREMENT, 
        EMAIL TEXT, 
        USERNAME TEXT, 
        PASSWORDHASH TEXT, 
        TOKEN TEXT, 
        IMG TEXT DEFAULT null,
        BG TEXT DEFAULT null,
        BIO TEXT DEFAULT null,
        ADDR TEXT DEFAULT null
    );

    CREATE TABLE POSTS (
        ID INTEGER PRIMARY KEY AUTOINCREMENT,
        USER_ID INTEGER,
        Text TEXT,
        IMG TEXT
    );

    CREATE TABLE COMMENTS (
        ID INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid INTEGER,
        post_id integer,
        comment_text TEXT
    );

    CREATE TABLE LIKES (
        ID INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid INTEGER,
        post_id INTEGER
    );
    
```



### Files and folders

- `src` The folder that holds my whole project.
- [`src\Structures.go`](https://github.com/Moody0101-X/Go_Api/blob/main/src/Structures.go) the file that contains all the models to be used in parsing and encapsulating data exmp => User, Post, Response, LoginForm....
- [`src\main.go`](https://github.com/Moody0101-X/Go_Api/blob/main/src/main.go) The program entry point, it run the server.
- [`src\routes.go`](https://github.com/Moody0101-X/Go_Api/blob/main/src/routes.go) contains routing functions that take the *gin.context* and handles the requests whether it is a POST or a GET request.
- [`src\cryptography.go`](https://github.com/Moody0101-X/Go_Api/blob/main/src/cryptography.go) contains utility functions for crypto operations like:
    1. Decode/Encode JWT.
    2. Hash passwords
    3. generate secret token for users.
- [`src\database.go`](https://github.com/Moody0101-X/Go_Api/blob/main/src/database.go) contains database functionality, given a global *sql.db* object to perform sqlite query.

### CDN

- I have recently added a cdn to be connected to once the api has images and other file to save and retrieve
Here is the [link](https://github.com/Moody0101-X/Zimg_cdn)
#### CDN connection

```go

// Go Constant endpoints
// server running in port 8500 so 8888 -> 8500?
const api string = "http://localhost:8500"
const addIMG string = api + "/Zimg/addAvatar"
const addBG string = api + "/Zimg/addbg"
const addPOST string = api + "/Zimg/NewPostImg"

```
    - `/Zimg/addAvatar` adding user avatar.
    - `/Zimg/addbg` adding user background.
    - `/Zimg/NewPostImg` adding user post assets. (Now only images. video will be added.)
    - To see implementation [cdn.go](https://github.com/Moody0101-X/Go_Api/blob/main/src/cdn.go)

### front-end app.

- to see the front-end app that is using this api go [Here](https://github.com/Moody0101-X/SM_app)
