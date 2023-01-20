package models;

type Comment struct {
	Id_          int `json:"id_"`
	Post_id		 int `json:"post_id"`
	Uuid		 int `json:"uuid"`
	Text		 string `json:"text"`
	User_		 AUser `json:"user"` // Filled when fitching comments.
}

type Like struct {
	Id_          int `json:"id_"`
	Post_id		 int `json:"post_id"`
	Uuid		 int `json:"uuid"`
	User_		 AUser `json:"user"` // Filled when fitching comments.
}

type TFollow struct {
	Id_        		int `json:"id_"`
	Follower_id		int `json:"follower_id"`
	Followed_id		int `json:"followed_id"`
	UToken			string `json:"token"`
}

type Follow struct {
	Id_        		int `json:"id_"`
	Follower_id		int `json:"follower_id"`
	Followed_id		int `json:"follower_id"`
}

type TokenizedComment struct {
	Post_id		 int `json:"post_id"`
	Uuid		 int `json:"uuid"`
	Text		 string `json:"text"`
	Token        string `json:"token"`
	Post_owner_id	 int `json:"post_owner_id"`
}

type TokenizedLike struct {
	Post_id		 int `json:"post_id"`
	Uuid		 int `json:"uuid"`
	Token        string `json:"token"`
	Post_owner_id	 int `json:"post_owner_id"`
}

type UserLogin struct {
	Password string `json:"Password"`
	Email    string `json:"Email"`
	Token 	 string `json:"T"`
}
