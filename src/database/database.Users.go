package database;

import (
	"fmt"
	"github.com/Moody0101-X/Go_Api/models"
	"github.com/Moody0101-X/Go_Api/crypto"
	"github.com/Moody0101-X/Go_Api/cdn"
)

func AuthenticateUserJWT(UserJWT string) models.Response {
    Token, Ok := crypto.GetTokenFromJwt(UserJWT);
    
    if Ok {
     
        User_, err := GetUserByToken(Token)
        if err != nil {
            // a db error.
            return models.MakeServerResponse(500, "Db Error. (line 108).")
        } else {
            // Returns the user if everything was alright.
            return models.MakeServerResponse(200, User_)
        }

    } else {
        // JWT error
        return models.MakeServerResponse(500, "server could not decode the token. (line 117)")
    }
}

func AuthenticateUserByEmailAndPwd(Pwd string, Email string) (models.User, models.Result) {
	var EmptyUser models.User

	if CheckUser(Email) {

		var user models.User
		row, err := dataBase.Query("SELECT PASSWORDHASH FROM USERS WHERE EMAIL=?", Email)
		
		defer row.Close()

		if err != nil {
			fmt.Println(err)
			return user, models.MakeServerResult(false, "Could not get user from db. 82")
		}

		var pwdHash string

		for row.Next() {
			row.Scan(&pwdHash)
		}
		
		if crypto.Sha256_(Pwd) == pwdHash {
			row, err := dataBase.Query("SELECT ID, EMAIL, USERNAME, TOKEN, IMG, BG, BIO, ADDR FROM USERS WHERE EMAIL=? ORDER BY ID DESC", Email)

			defer row.Close()

			if err != nil {
				return user, models.MakeServerResult(false, "Could not get user from db. 97")
			}

			for row.Next() {
				row.Scan(&user.Id_, &user.Email, &user.UserName, &user.Token, &user.Img, &user.Bg,  &user.Bio, &user.Address)
				user.Img = CheckCdnLink(user.Img);
				user.Bg = CheckCdnLink(user.Bg);
			}

			JWT, err := crypto.StoreTokenInJWT(user.Token)

			if err == nil {
				user.Token = JWT
				return user, models.MakeServerResult(true, "User created! you can login now..")
			}


			return EmptyUser, models.MakeServerResult(false, "Server had a problem encoding the token..")
		}
		
		return user, models.MakeServerResult(false, "incorrect password. try again")
	}

	return EmptyUser, models.MakeServerResult(false, "incorrect Email.. check and try again!")
}

func GetUserById(id int) models.AUser {
	
	var User models.AUser

	row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE ID=? ORDER BY ID DESC", id)
	defer row.Close()
	
	if err != nil {
		fmt.Println(err)
		return User
	}

	for row.Next() {
		row.Scan(&User.Id_, &User.UserName,&User.Img, &User.Bg, &User.Bio, &User.Address)
		User.Img = CheckCdnLink(User.Img);
		User.Bg = CheckCdnLink(User.Bg);
	}

	return User
}

func GetUserByToken(Token string) (models.User, error) {
	
	var User_ models.User
	row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE TOKEN=? ORDER BY ID DESC", Token)
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return User_, err
	}

	for row.Next() {
		row.Scan(&User_.Id_, &User_.UserName, &User_.Img, &User_.Bg, &User_.Bio, &User_.Address)
		User_.Img = CheckCdnLink(User_.Img);
		User_.Bg = CheckCdnLink(User_.Bg);
	}

	return User_, nil
}

func GetUsers(uuid interface{}) []models.AUser {
	var Users []models.AUser
	switch uuid.(type) {
		case int:
			row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS ORDER BY ID DESC")
			defer row.Close()
			
			if err != nil {
				fmt.Println(err)
				return Users
			}

			var temp models.AUser

			for row.Next() {
				row.Scan(&temp.Id_, &temp.UserName, &temp.Img, &temp.Bg, &temp.Bio, &temp.Address)
				temp.IsFollowed = IsFollowing(temp.Id_, uuid.(int))
				temp.Img = CheckCdnLink(temp.Img);
				temp.Bg = CheckCdnLink(temp.Bg);
				Users = append(Users, temp)
			}

			break

		case bool:

			row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS ORDER BY ID DESC")
			defer row.Close()
			
			if err != nil {
				fmt.Println(err)
				return Users
			}

			var temp models.AUser

			for row.Next() {
				
				row.Scan(&temp.Id_, &temp.UserName, &temp.Img, &temp.Bg, &temp.Bio, &temp.Address)
				temp.Img = CheckCdnLink(temp.Img);
				temp.Bg = CheckCdnLink(temp.Bg);
				Users = append(Users, temp)
			}

			break

		default:
			return Users
			break
	}
	

	return Users	
}

func GetUsersByQuery(Q string, uuid interface{}) []models.AUser {
	var Users []models.AUser
	var NewQ string = "%" + Q + "%"

	switch uuid.(type) {
		case int:
			row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE USERNAME LIKE ? ORDER BY ID DESC", NewQ)

			defer row.Close()

			if err != nil {
				fmt.Println(err)
				return Users
			}

			for row.Next() {
				var temp models.AUser
				row.Scan(&temp.Id_, &temp.UserName, &temp.Img, &temp.Bg, &temp.Bio, &temp.Address)
				temp.IsFollowed = IsFollowing(temp.Id_, int(uuid.(int)))
				temp.Img = CheckCdnLink(temp.Img);
				temp.Bg = CheckCdnLink(temp.Bg);
				Users = append(Users, temp)	
			}

			break

		case bool:
			
			row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE USERNAME LIKE ? ORDER BY ID DESC", NewQ)

			defer row.Close()

			if err != nil {
				fmt.Println(err)
				return Users
			}

			for row.Next() {
				var temp models.AUser
				row.Scan(&temp.Id_, &temp.UserName, &temp.Img, &temp.Bg, &temp.Bio, &temp.Address)
				temp.Img = CheckCdnLink(temp.Img);
				temp.Bg = CheckCdnLink(temp.Bg);
				Users = append(Users, temp)	
			}

			break

		default:
			return Users		
	}

	return Users
}

func GetUserByJWToken(JWToken string) (models.User, bool) {
	Token, isValid := crypto.GetTokenFromJwt(JWToken);

	if isValid {
		User, _ := GetUserByToken(Token)
		return User, true
	} else {
		var EmptyUser models.User
		return EmptyUser, true
	}
}

func GetUserIdByToken(t string) (int, bool) {
	var id int
	row, err := dataBase.Query("SELECT ID FROM USERS WHERE TOKEN=?", t)
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return id, false
	}

	for row.Next() {
		row.Scan(&id)
	}

	return id, true
}

func CheckUser(Email string) bool {
		
	row, err := dataBase.Query("SELECT ID FROM USERS WHERE EMAIL=? ORDER BY ID DESC", Email)
	defer row.Close()
	var u []int;

	if err != nil {
		fmt.Println(err)
		return false
	}

	var id int;

	for row.Next() {	
		row.Scan(&id)
		u = append(u, id)
	}

	return (len(u) >= 1);
}

func AddUser(user models.User) models.Response {
	if !CheckUser(user.Email) {

		var Token string = crypto.GenerateAccessToken(user.Email)
		

		/*------------Add To cdn-------------*/
		var uuid = GetNextUID("Users")

		ok, img := cdn.AddUserAvatarToCdn(uuid, user.Img)

		if !ok {
			return models.MakeServerResponse(500, "cdn error, could not add avatar.")
		}

		ok, bg := cdn.AddUserBackgroundToCdn(uuid, user.Bg)

		if !ok {
			return models.MakeServerResponse(500, "cdn error, could not add background.")
		}


		/*------------Add To cdn-------------*/

		stmt, _ := dataBase.Prepare("INSERT INTO USERS(EMAIL, USERNAME, PASSWORDHASH, TOKEN, IMG, BIO, BG, ADDR) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
		_, err := stmt.Exec(user.Email, user.UserName, user.PasswordHash, Token, img, user.Bio, bg, user.Address)
		
		if err != nil {
			return models.MakeServerResponse(500, "Could not add to db.")
		}

		FetchedUser, err := GetUserByToken(Token)

		if err != nil {
			return models.MakeServerResponse(500, "Could not get created user from db. L288")
		} else {
			
			JWT, err := crypto.StoreTokenInJWT(Token)
			
			if err != nil {
				fmt.Println(err)
				return models.MakeServerResponse(500, "The server had a problem making the jwt token.")
			}

			FetchedUser.Token = JWT;
			return models.MakeServerResponse(200, FetchedUser) // Success.
		}

	}

	return models.MakeServerResponse(500, "This user already exists..")
}

func GetuuidByToken(Token string) (int, bool) {
	
	var uuid int
	row, err := dataBase.Query("SELECT ID FROM USERS WHERE TOKEN=?", Token)
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return 0, false
	}

	for row.Next() {
		row.Scan(&uuid);
	}

	return uuid, true
}

func UpdateUser(field string, newValue string, Token string) models.Result {
	
	var ok bool;
	var Query string
	
	switch field {
		case "IMG":
			Query = "UPDATE USERS SET IMG=? WHERE TOKEN=?"
			uuid, OK := GetuuidByToken(Token)
			
			if OK {
				OK, newValue = cdn.AddUserAvatarToCdn(uuid, newValue)
				fmt.Println("path: ", newValue)
				ok = OK
			}

			break		

		case "BIO":
			Query = "UPDATE USERS SET BIO=? WHERE TOKEN=?"
			ok = true
			break

		case "ADDR":
			Query = "UPDATE USERS SET ADDR=? WHERE TOKEN=?"
			ok = true
			break		
		case "BG":
			
			Query = "UPDATE USERS SET BG=? WHERE TOKEN=?"
			
			uuid, OK := GetuuidByToken(Token)
			
			if OK {
				OK, newValue = cdn.AddUserBackgroundToCdn(uuid, newValue)
				ok = OK
			}

			break
		
		case "USERNAME":
			Query = "UPDATE USERS SET USERNAME=? WHERE TOKEN=?"
			ok = true
			break

		default:
			ok = false
			break
	}

	if ok {

		stmt, _ := dataBase.Prepare(Query)
		_, err := stmt.Exec(newValue, Token)

		if err != nil {
			fmt.Println("db err: ", err)
			return models.MakeServerResult(false, "db err, could not update.")
		}

		return models.MakeServerResult(true, "success!")
	}

	return models.MakeServerResult(false, "Unexpected field name, or could not find user by token..")

}
