package entity

type User struct {
	Name      string `db:"name" json:"name"`
	Img       string `db:"img" json:"img"`
	Email     string `db:"email" json:"email"`
	Dob       string `db:"dob" json:"dob"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type UpdateUser struct {
	Name     string `form:"name" json:"name"`
	Img      string `form:"img" json:"img"`
	Email    string `form:"email" json:"email"`
	Dob      string `form:"dob" json:"dob"`
	Password string `form:"password" json:"password"`
}
