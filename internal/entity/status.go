package entity

type Status struct {
	IsLogin  bool `json:"is_login"`
	IsSeller bool `json:"is_seller"`
	IsAdmin  bool `json:"is_admin"`
}
