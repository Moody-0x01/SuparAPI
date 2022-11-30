

# Go Api

### Routes

NOTE: This section will be updated regularily once I add or remove a route.
- `/login` Login, can get a token or (password, Email) and if everything is valid it returns user info + Token.
- `/getUserPosts` return the user's posts, and it gets an id_ form value as uuid.
- `/GetAllPosts` return all the posts from db.
- `/query` used to query a specific user. if not found then an empty list is returned or an error response code.
- `/:uuid` quick user lookup by uuid.
- `/signup` making an account. then returning the token to update, delete, add and other operation regarding ur account.

### Files and folders

- `src` The folder that holds my whole project.
- `src\structures.go` the file that contains all the models to be used in parsing and encapsulating data exmp => User, Post, Response, LoginForm....
- `src\main.go` The program entry point, it run the server.
- `src\WebApi.go` contains routing functions that take the *gin.context* and handles the requests whether it is a POST or a GET request.
- `src\ApplicationApi.go` contains utility functions to do this ops:
    1. Decode/Encode JWT.
    2. Hash passwords
    3. generate secret token for users.
- `src\Database_functions` contains database functionality, given a global *sql.db* object to perform sqlite query.
- `src\Db_setup.go` Initializes the database global connection. if there is an error connecting to the db, the api would not work.




















































