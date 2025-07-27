package entity

type SignUp struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Dob      string `json:"dob"`
	Password string `json:"password"`
}
