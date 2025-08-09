package entity

type User struct {
	ID        string `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Img       string `db:"img" json:"img"`
	Email     string `db:"email" json:"email"`
	Dob       string `db:"dob" json:"dob"`
	Password  string `db:"password" json:"password"`
	CreatedAt string `db:"created_at" json:"created_at"`
}
