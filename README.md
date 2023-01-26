
# Go Api
### Routes
NOTE: This section will be updated regularily once I add or remove a route.
NOTE: alot of the front-end routes are locked if you are not logged in...

- `/` the home page for my app. also can be accessed from `/Home` but only if you are logged in because you need a user profile to see posts and like/comment under them. it also has a post feature. img or text or both..s
- `/signup` the page in which you can make an account
- `/login` the page in which you can login
- `/Accounts` the page to visualize all user account and look for someone.
- `/Accounts/:uuid` see profile of a user by uuid.
- `/v2/login` Login, can get a token or (password, Email) and if everything is valid it returns user info + Token.
- `/v2/getUserPosts` return the user's posts, and it gets an id_ form value as uuid.
- `/v2/GetAllPosts` return all the posts from db.
- `/v2/query` used to query a specific user. if not found then an empty list is returned or an error response code.
- `/v2/:uuid` quick user lookup by uuid.
- `/v2/signup` making an account. then returning the token to update, delete, add and other operation regarding ur account.
- `/v2/update` expects either a img, bio, addr or bg to be updated in the DATABASE, but also accepts a token that will be given if aut/v2henticated. if the token is not provided, the server will return error code 500.
- `/v2/NewPost` add new post with token, expects a token, img, uuid, and post text. if something is not set properly it will return a response report about the error.
- `/v2/DeletePost` deleting user post using userid, postid and jwt.
- `/v2/comment` a route for commenting in posts and stuff.
- `/v2/like` a route to add likes to posts
- `/v2/like/remove` a route to remove a like from a posts.
- `/v2/follow` a follow endpoint to follow users.
- `/v2/unfollow` a route for unfollowing users.
- `/v2/getFollowers/:uuid` getting the followers of a user. by uuid..
- `/v2/getFollowings/:uuid` getting the followings of a user. by uuid..
- `/v2/getComments/:pid` a route to get comments by pid (post id).
- `/v2/getLikes/:pid` a route to get the likes of a certain post by id.
- `/v2/chat/getUserConversations` a route to a user's conversations data, as an array (which will be changed due to some technical reasons)
- `/v2/chat/getUserConversationById` a route that return one conversation by its id
- `/v2/WSoc` a route that handles tcp connections via a socket that connects from `ws://host:port/v2/WSoc`

NOTE: an interesting update is that all endpoints are prefixed by /v2 because I added my single page app and integrated the front-end route so anything that does not have /v2 is a front-end thing that returns html. and the opposite is an api endpoint.

### DATABASE (Tables and schema.)

I have used only 2 tables just because I did not want to overcomplicate things, but more will be added as m progressing in this project.

- Check [`src\script.sql`](https://github.com/Moody0101-X/Go_Api/blob/main/src\script.sql)
- Note: Adding more data fields and appropriate type is kinda crucial, but this is it for now.

### Files and folders

- `src` The folder that holds my whole project.
- [`src\models`](https://github.com/Moody0101-X/Go_Api/blob/main/src/models) the dirs that contains all the models to be used in parsing and encapsulating data exmp => User, Post, Response, LoginForm....

- [`src\main.go`](https://github.com/Moody0101-X/Go_Api/blob/main/src/main.go) The program entry point, it run the server.
- [`src\routes`](https://github.com/Moody0101-X/Go_Api/blob/main/src/routes) contains routing functions that take the *gin.context* and handles the requests whether it is a POST or a GET request.
- [`src\crypto`](https://github.com/Moody0101-X/Go_Api/blob/main/src/crypto) contains utility functions for crypto operations like:
    
    1. Decode/Encode JWT.
    2. Hash passwords
    3. generate secret token for users.

- [`src\database`](https://github.com/Moody0101-X/Go_Api/blob/main/src/database) contains DATABASE functionality, given a global *sql.db* object to perform sqlite queries, and I recently added a socket connection support for real time data transmition! (Updates, messages).


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

### Issues.

## Front-end
    - Messy code all over the app.
    - Non-standarized Ui main components.
    - Responsible for parsing the data once it comes in, which is not that convinient I guess, we need to do everything in the backend then return the needed result so the ui has one last thing to do, change state.
    - Refactoring needed
    
## Back-end
    - Response is not standard, we need new data struct to return a parsed and more organized data for big data objects
    
    ### Example
```javascript
    // Instead of this ob.
    const c = [
            {
                id_:0,
                msgs: [{...}],
                ts:new Date(),
            }, {
                id_:1,
                msgs: [{...}],
                msgCount: msgs.length,
                ts:new Date(),
            }
        ]


    // We need this.
    const c1 = {
        0: {
            msgs: [{...}],
            msgCount: msgs.length,
            ts:new Date(),
        },
        1: {
            msgs: [{...}],
            msgCount: msgs.length,
            ts:new Date(),
        }
    }

    // The purpose is obvs so we just prevent this code
    c.map((conv) => {
        
        if(conv.id_ === target) {
            // DO something if the targeted conv is found.    
        }
    })

    // Istead.
    const f = (conversation) => {
        // DO something if the targeted conv is found.
    }
    const targetedC = c[target]
    // Then operate.. change state or append, do whatever but it is just faster, way faster than maping.
    f(targetedC)

    /*
    
    But idk, it just makes more sense, because I am already using maps all over the place..
    
    */

```
    - Refactoring needed

## CDN
    - More file formats (mp3, mp4, mkv-maybe....).
    - needs more effecient design.
    - Refactoring needed