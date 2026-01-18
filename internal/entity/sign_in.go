package entity

type SignIn struct {
	Id       string `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type UserPayment struct {
	Name  string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
}
