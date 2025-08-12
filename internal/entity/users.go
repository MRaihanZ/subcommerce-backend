package entity

import "time"

type User struct {
	Name      string `db:"name" json:"name"`
	Img       string `db:"img" json:"img"`
	Email     string `db:"email" json:"email"`
	Dob       string `db:"dob" json:"dob"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type UpdateUser struct {
	Name     string    `form:"name"`
	ImgPath  string    `form:"imgPath"`
	Email    string    `form:"email"`
	Dob      time.Time `form:"dob" time_format:"2006-01-02"`
	Password string    `form:"password"`
}
