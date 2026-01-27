package entity

type Admin struct {
	Id    string `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
}

type CreateAdmin struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

type UpdateAdmin struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

type PayoutEmailData struct {
	SellerName string `db:"name"`
	Email      string `db:"email"`
	Name       string `db:"transfer_name"`
	Amount     int64  `db:"transfer_amount"`
}
