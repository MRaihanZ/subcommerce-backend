package entity

type Response[D any] struct {
	Code   int     `json:"code"`
	Status string  `json:"status"`
	Data   D       `json:"data"`
	Error  *string `json:"error"`
}
