package models;

type User struct {
	Id_ 		 int `json:"id_"`
	Email 		 string `json:"Email"`
	UserName 	 string `json:"UserName"`
	PasswordHash string `json:"Password"`
	Token 		 string `json:"token"`
	Img 		 string `json:"img"`
	Bg 			 string `json:"bg"`
	Bio 		 string `json:"bio"`
	Address		 string `json:"addr"`
}

type AUser struct {
	Id_ 		 int `json:"id_"`
	UserName 	 string `json:"UserName"`
	Img		 	 string `json:"img"`
	Bg 			 string `json:"bg"`
	Bio 		 string `json:"bio"`
	Address		 string `json:"addr"`
	IsFollowed	 bool `json:"isfollowed"`
}

func (U *User) SetDefaults() {
	//TODO Setting the default fields to add to the db if some are not present.

	if isEmpty(U.Img) {
		U.Img = DefaultUserImg
	}
	
	if isEmpty(U.Bg) {
		U.Bg = DefaultUserBg
	}
	
	if isEmpty(U.Bio) {
		U.Bio = DefaultUserBio
	}

	if isEmpty(U.Address) {
		U.Address = DefaultUserAddress
	}
}
