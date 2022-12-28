package models;


type Notification struct {
	Id_ int	`json:"id"`
	Text string	`json:"text"`
	Type string	`json:"type"` 
	Uuid int	`json:"uuid"`
	Actorid int	`json:"actorid"`
	Seen bool	`json:"seen"`
	Post_id int	`json:"post_id"`
	Link  string `json:"link"`
}

func NewNot(
	id int, text string, t int, uuid int, actorid int, Post_id int
) Notification {
	var new Notification;
	
	new.Id_ = id;
	new.Text = text;
	new.Type = t;
	new.Uuid = uuid;
	new.Actorid = actorid;
	new.Seen = 0;
	new.Post_id = Post_id;
	new.Link = "";
	
	return new;
}

