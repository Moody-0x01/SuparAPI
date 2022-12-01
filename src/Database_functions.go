package main;

import "fmt"

func db() {
	fmt.Println("Hi thhis is db!");
}

// TODO verifyUserPassword, AddUser, MakePost.

func print(s string) {
	fmt.Println(s);
}


/*-------------------------------------------------------------------------------------------------------------------------------
 	POSTS
-------------------------------------------------------------------------------------------------------------------------------*/


// func verifyUserPassword(LoginObject UserLogin) (User, bool) {
// 	// TODO	
// }


// func MakePost(post Post) bool {
// 	//TODO
// }

func GetAllPosts() []Post {
	var Posts []Post

	row, err := dataBase.Query("SELECT Text, IMG FROM POSTS ORDER BY ID DESC")
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return Posts
	}

	var temp Post

	for row.Next() {
		row.Scan(&temp.Text, &temp.Img);
		Posts = append(Posts, temp);
	}

	fmt.Println(Posts);
	return Posts
}

func getUserPostById(id string) []Post {
	// A functions to use 
	var Posts []Post

	row, err := dataBase.Query("SELECT Text, IMG FROM POSTS WHERE USER_ID=? ORDER BY ID DESC", id)
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return Posts
	}

	var temp Post

	for row.Next() {
		row.Scan(&temp.Text, &temp.Img);
		Posts = append(Posts, temp);
	}

	return Posts
}



func AuthenticateUserByEmailAndPwd(Pwd string, Email string) (User, Error) {
	
	var user User

	// row, err := dataBase.Query("SELECT ID, EMAIL, PASSWORDHASH, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE ID=? ORDER BY ID DESC", id)
	row, err := dataBase.Query("SELECT PASSWORDHASH FROM USERS WHERE EMAIL=?", Email)
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return user, MakeServerError(false, "Could not get user from db. 82")
	}

	var pwdHash string

	for row.Next() {
		row.Scan(&pwdHash)
	}
	
	if sha256_(Pwd) == pwdHash {
		//TODO Get user.
		//TODO Return user.
		row, err := dataBase.Query("SELECT ID, EMAIL, USERNAME, TOKEN, IMG, BG, BIO, ADDR FROM USERS WHERE EMAIL=? ORDER BY ID DESC", Email)

		defer row.Close()

		if err != nil {
			return user, MakeServerError(false, "Could not get user from db. 97")
		}

		for row.Next() {
			row.Scan(&user.Id_, &user.Email, &user.UserName, &user.Token, &user.Img, &user.Bg,  &user.Bio, &user.Address)
		}

		JWT, err := StoreTokenInToken(user.Token)

		if err == nil {
			user.Token = JWT
			return user, MakeServerError(true, "User created! you can login now..")
		}

		var EmptyUser User
		return EmptyUser, MakeServerError(false, "Server had a problem encoding the token..")

	}
	
	return user, MakeServerError(false, "incorrect password. try again")
}

/*-------------------------------------------------------------------------------------------------------------------------------
 	USERS
-------------------------------------------------------------------------------------------------------------------------------*/
func getUserById(id string) User {
	
	var User User

	row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE ID=? ORDER BY ID DESC", id)
	defer row.Close()
	if err != nil {
		fmt.Println(err)
		return User
	}

	for row.Next() {
		row.Scan(&User.Id_, &User.UserName,&User.Img, &User.Bg, &User.Bio, &User.Address)
	}



	return User
}

func getUserByToken(Token string) (User, error) {
	var User User

	row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE TOKEN=? ORDER BY ID DESC", Token)
	
	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return User, err
	}

	for row.Next() {
		row.Scan(&User.Id_, &User.UserName, &User.Img, &User.Bg, &User.Bio, &User.Address)
	}

	return User, nil
}

func getUsers() []User {
	var Users []User
	row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS ORDER BY ID DESC")
	defer row.Close()
	if err != nil {
		fmt.Println(err)
		return Users
	}

	var temp User

	for row.Next() {
		
		row.Scan(&temp.Id_, &temp.UserName, &temp.Img, &temp.Bg, &temp.Bio, &temp.Address)
		Users = append(Users, temp)
	}

	return Users	
}

func getUsersByQuery(Q string) []User {
	// TODO Verify this one works, sqlite to find alike attributes in row.
	// DONE.
	var Users []User
	var NewQ string = "%" + Q + "%"
	row, err := dataBase.Query("SELECT ID, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE USERNAME LIKE ? ORDER BY ID DESC", NewQ)

	defer row.Close()

	if err != nil {
		fmt.Println(err)
		return Users
	}

	for row.Next() {
		var temp User
		row.Scan(&temp.Id_, &temp.UserName, &temp.Img, &temp.Bg, &temp.Bio, &temp.Address)
		Users = append(Users, temp)	
	}

	return Users
}

func GetUserByJWToken(JWToken string) (User, bool) {
	Token, isValid := GetTokenFromJwt(JWToken);

	if isValid {
		return GetUserByToken(Token), true
	} else {
		var EmptyUser User
		return EmptyUser, true
	}
}

func GetUserByToken(Token string) User {
	var User User;
	
	row, err := dataBase.Query("SELECT ID, EMAIL, USERNAME, IMG, BG, BIO, ADDR FROM USERS WHERE TOKEN=? ORDER BY ID DESC", Token)
	defer row.Close()
	if err != nil {
		fmt.Println(err)
		return User
	}

	for row.Next() {
		row.Scan(&User.Id_, &User.Email, &User.UserName, &User.Img, &User.Bg, &User.Bio, &User.Address)
	}

	return User;
}


func CheckUser(Email string) bool {
		
	row, err := dataBase.Query("SELECT ID FROM USERS WHERE EMAIL=? ORDER BY ID DESC", Email)
	defer row.Close()
	var u []User;

	if err != nil {
		fmt.Println(err)
		return false
	}

	var temp User

	for row.Next() {	
		row.Scan(&temp.Id_, &temp.UserName, &temp.Img, &temp.Bg, &temp.Bio, &temp.Address)
		u = append(u, temp)	
	}

	if len(u) == 1 {
		return true
	} else if len(u) == 0 {
		return false
	} else {
		return true
	}
}

func AddUser(user User) Response {
	
	/* 
	- Check if the email associated with the request already exists.
		if it exists return an error code and a string.
		else just go on with process
	- try to add the user by:
		generating a token and assigning it to the user before adding
	*/

	

	if !CheckUser(user.Email) {

		var Token string = generateAccessToken(user.Email)
		stmt, _ := dataBase.Prepare("INSERT INTO USERS(EMAIL, USERNAME, PASSWORDHASH, TOKEN, IMG, BIO, BG, ADDR) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
		_, err := stmt.Exec(user.Email, user.UserName, user.PasswordHash, Token, user.Img, user.Bio, user.Bg, user.Address)
		
		if err != nil {
			return MakeServerResponse(500, "Could not add to db.")
		}

		FetchedUser, err := getUserByToken(Token)

		if err != nil {
			return MakeServerResponse(500, "Could not get created user from db. L288")
		} else {
			
			JWT, err := StoreTokenInToken(Token)
			
			if err != nil {
				fmt.Println(err)
				return MakeServerResponse(500, "The server had a problem making the jwt token.")
			}

			fmt.Println("JWT: ", JWT)

			FetchedUser.Token = JWT;
			return MakeServerResponse(200, FetchedUser) // Success.
		}

	}

	return MakeServerResponse(500, "This user already exists..")
}


func updateUser(field string, newValue string, Token string) Error {
	
	stmt, _ := dataBase.Prepare("UPDATE USERS SET ?=? WHERE TOKEN=?")
	_, err := stmt.Exec(field, newValue, Token)

	if err != nil {
		return MakeServerError(false, "db err, could not update.")
	}

	return MakeServerError(true, "success!")
}


func AddPost(Text string, Img string, uuid int) Error {

	stmt, _ := dataBase.Prepare("INSERT INTO POSTS(USER_ID, Text, IMG) VALUES(?, ?, ?)")
	_, err := stmt.Exec(uuid, Text, Img)
	if err != nil {
		return MakeServerError(false, "could not add post. err L334")
	}

	return MakeServerError(true, "success!")
	
}










// /*-------------------------------------------------------------------------------------------------------------------------------
// Get Specific rows with a condition.
// ---------------------------------------------------------------------------------------------------------------------------------------*/
// // GET ROW FROM {tableName} WHERE {fieldName} == {value}

// func getItemByValue(c *gin.Context, tableName string, fieldName string, Value string) {
// 	// Gets the product by category
// 	var QueryString string
// 	if tableName == "PRODUCTS" {
// 		QueryString = "SELECT ID, IMG, NAME, DESC, CATEGORY FROM " + tableName + " WHERE " + fieldName  + " = " + Value
// 	} else {
// 		QueryString = "SELECT * FROM " + tableName + " WHERE " + fieldName + " = " + Value
// 	}

// 	print(QueryString)
// 	row, err := dataBase.Query(QueryString)

// 	if err != nil {
// 		c.JSON(http.StatusOK, newResp("Erreur Base de donne"))
// 		return
// 	}

// 	defer row.Close()

// 	resp := newResp("ok")
// 	switch(tableName) {
// 		case "PRODUCTS":
// 			var id_, img, NAME, Desc, category string

// 			for row.Next() {
// 				n := newHzJson()
// 				row.Scan(&id_, &img, &NAME, &Desc, &category)
// 				n.setAttribute("ID", id_)
// 				n.setAttribute("img", img)
// 				n.setAttribute("Name", NAME)
// 				n.setAttribute("category", category)
// 				resp.appendChild(n)
// 			}

// 			break
// 		case "CONTACT":
			
// 			var FullName, Email, Message string
// 			for row.Next() {
// 				n := newHzJson()
// 				row.Scan(&FullName, &Email, &Message)
				
// 				fmt.Println("FullName: ", FullName, " Email: ", Email, " Message: ", Message)

// 				n.setAttribute("FullName", FullName)
// 				n.setAttribute("Email", Email)
// 				n.setAttribute("Name", Message)
// 				resp.appendChild(n)
// 			}

// 			break
// 		case "ORDERS":
// 			var FullName, PhoneNumber, Address, City, ProductID string

// 			for row.Next() {
// 				n := newHzJson()
// 				row.Scan(&FullName, &PhoneNumber, &Address, &City, &ProductID)
				
// 				fmt.Println("FullName: ", FullName, "PhoneNumber: ", PhoneNumber, "Address: ", Address,"City: ", City, "ProductID", ProductID)

// 				n.setAttribute("FullName", FullName)
// 				n.setAttribute("PhoneNumber", PhoneNumber)
// 				n.setAttribute("Address", Address)
// 				n.setAttribute("City", City)
// 				n.setAttribute("ProductID", ProductID)
// 				resp.appendChild(n)
// 			}

// 			break
// 		case "APPLICATIONS":
// 			// FULLNAME TEXT,
// 			// EMAIL TEXT,
// 			// TELE TEXT,
// 			// CITY TEXT,
// 			// LINKDIN TEXT,
// 			// CV TEXT
// 			var FullName, Email, PhoneNumber, City, Linkdin, cvPath string

// 			for row.Next() {
// 				n := newHzJson()
// 				row.Scan(&FullName, &Email, &PhoneNumber, &City, &Linkdin, &cvPath)
				
// 				fmt.Println(FullName, "FullName", Email, "Email", PhoneNumber, "PhoneNumber",City, "City",Linkdin, "Linkdin",cvPath, "cvPat")

// 				n.setAttribute("FullName", FullName)
// 				n.setAttribute("Email", Email)
// 				n.setAttribute("PhoneNumber", PhoneNumber)
// 				n.setAttribute("City", City)
// 				n.setAttribute("Linkdin", Linkdin)
// 				n.setAttribute("cvPat", cvPath)
// 				resp.appendChild(n)
// 			}

// 			break
// 		default:
// 			c.JSON(http.StatusOK, newResp("Unsopprted method..."))
// 	}
	
// 	if mapB, err := json.Marshal(resp); err != nil {
// 			fmt.Println("error recv: ", err.Error())
// 	} else {
// 		fmt.Println("response:", string(mapB))
// 	}

// 	c.JSON(http.StatusOK, resp)
// }

// func getProductById(c *gin.Context, ID string) {
// 	getItemByValue(c, "PRODUCTS", "ID", ID)
// }

// func getProductsByCategory(c *gin.Context, category string) {
// 	getItemByValue(c, "PRODUCTS", "CATEGORY", category)
// }

// func getProductByMark(c *gin.Context, mark string) {
// 	getItemByValue(c, "PRODUCTS", "MARK", mark)
// }


// func getProductPrice(c *gin.Context, ID string) {
// 	return 
// }

// func getOrderByName(c *gin.Context, json *HzJson) {
// 	FullName, _ := json.getAttribute("Name");
// 	getItemByValue(c, "ORDERS", "FULLNAME", FullName)
// }


// func addApplication(c *gin.Context, json *HzJson) {
// 	/*
// 		Data:
// 			FullName,
// 			Email,
// 			PhoneNumber,
// 			City,
// 			Linkdin,
// 			position,
// 			cv.
// 	*/

// 	FullName, _ := json.getAttribute("FullName")
// 	Email, _ := json.getAttribute("Email")
// 	PhoneNumber, _ := json.getAttribute("PhoneNumber")
// 	City, _ := json.getAttribute("City")
// 	Linkdin, _ := json.getAttribute("Linkdin")
// 	position, _ := json.getAttribute("position")
// 	cv, _ := json.getAttribute("cv")


// 	fmt.Println(FullName)
// 	fmt.Println(Email)
// 	fmt.Println(PhoneNumber)
// 	fmt.Println(City)
// 	fmt.Println(Linkdin)
// 	fmt.Println(position)
// 	fmt.Println(cv)




// 	// _, err := dataBase.Query("INSERT INTO APPLICATIONS (FULLNAME, EMAIL, TELE, CITY, LINKDIN, CV) values (?, ?, ?, ?, ?, ?)", FullName, Email, PhoneNumber, City, Linkdin, position, cv)

// 	// if err != nil {
// 	// 	c.JSON(http.StatusOK, newResp("was not added"))
// 	// 	fmt.Print(err)
// 	// 	return
// 	// }
	
// 	c.JSON(http.StatusOK, newResp("ok"))
// }

// func addOrder(c *gin.Context, json *HzJson) {
// 	/*
// 		Data:
// 			FullName,
// 			Email,
// 			PhoneNumber,
// 			City,
// 			Linkdin,
// 			position,
// 			cv.
// 	*/

// 	FullName, _ := json.getAttribute("FullName")
// 	PhoneNumber, _ := json.getAttribute("PhoneNumber")
// 	City, _ := json.getAttribute("City")
// 	Address, _ := json.getAttribute("Address")
// 	_, err := dataBase.Query("INSERT INTO ORDERS (FULLNAME, PHONENUMBER, CITY, ADDRESS, PRODUCTS_ID) values (?, ?, ?, ?, ?)", FullName, PhoneNumber, City, Address)

// 	if err != nil {
// 		c.JSON(http.StatusOK, newResp("was not added"))
// 		fmt.Print(err)
// 		return
// 	}
	
// 	c.JSON(http.StatusOK, newResp("ok"))
// }

// func getLatestProducts(c *gin.Context, n int) {
// 	// gets the last n elements
// 	var MAX string;
// 	var newMAX int;

// 	row, err := dataBase.Query("SELECT max(CAST(ID AS INTEGER)) FROM PRODUCTS")
	
// 	if err != nil {
// 		print("Troble")
// 	}

// 	for row.Next() {
// 		row.Scan(&MAX)
// 		newMAX, _ = strconv.Atoi(MAX)
// 	}

// 	var newId string = strconv.Itoa(newMAX - n)

// 	prodrow, err := dataBase.Query("SELECT ID, IMG, NAME, DESC, CATEGORY FROM PRODUCTS WHERE CAST(ID AS INTEGER) > ?", newId)

// 	if err != nil {
// 		c.JSON(http.StatusOK, newResp("Erreur Base de donne"))
// 		print("Arror")
// 		return
// 	}

// 	var id_, img, NAME, Desc, category string
// 	resp := newResp("ok")

// 	for prodrow.Next() {

// 		n := newHzJson()
// 		fmt.Print(prodrow);
// 		prodrow.Scan(&id_, &img, &NAME, &Desc, &category)
// 		fmt.Println("ID: ", id_, "IMG: ", img, "NAME: ", NAME, "DESC: ", Desc)
// 		n.setAttribute("ID", id_)
// 		n.setAttribute("img", img)
// 		n.setAttribute("Name", NAME)
// 		n.setAttribute("category", category)
// 		resp.appendChild(n)

// 	}

// 	if mapB, err := json.Marshal(resp); err != nil {
// 			fmt.Println("error recv: ", err.Error())
// 	} else {
// 		fmt.Println("response:", string(mapB))
// 	}

// 	c.JSON(http.StatusOK, resp)
// }



// // Costum database calls.
// func getProductDescById(c *gin.Context, id string) {

// 	row, err := dataBase.Query("SELECT DESC, REF, MARK FROM PRODUCTS WHERE ID = ?", id)
	
// 	if err != nil {
// 		c.JSON(http.StatusOK, newResp("Erreur Base de donne"))
// 		fmt.Println("Arror")
// 		return
// 	}

// 	var Desc, Marque, Ref string
// 	resp := newResp("ok")

// 	for row.Next() {
// 		n := newHzJson()
// 		row.Scan(&Desc, &Ref, &Marque)

// 		fmt.Println("Id: ", id)
// 		fmt.Println("DESC: ", Desc)
// 		fmt.Println("MARK: ", Marque)
// 		fmt.Println("REF: " , Ref)
// 		n.setAttribute("DESC", Desc)
// 		n.setAttribute("REF", Ref)
// 		n.setAttribute("MARK", Marque)
// 		resp.appendChild(n)
// 	}

// 	if mapB, err := json.Marshal(resp); err != nil {
// 		fmt.Println("error recv: ", err.Error())
// 	} else {
// 		fmt.Println("response:", string(mapB))
// 	}

// 	c.JSON(http.StatusOK, resp)
// }



