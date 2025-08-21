package entity

type Rating struct {
	AverageRating float32 `db:"average_rating" json:"average_rating"`
	RatingCount   int     `db:"rating_count" json:"rating_count"`
}

type RatingComments struct {
	UName     string `db:"u_name" json:"u_name"`
	Img       string `db:"img" json:"img"`
	PvName    string `db:"pv_name" json:"pv_name"`
	Interval  int    `db:"interval" json:"interval"`
	IName     string `db:"i_name" json:"i_name"`
	Rating    int    `db:"rating" json:"rating"`
	Comment   string `db:"comment" json:"comment"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type RatingRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}
