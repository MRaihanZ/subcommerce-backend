package entities

type Response[D any] struct {
	Code   string  `json:"code"`
	Status string  `json:"status"`
	Data   D       `json:"data"`
	Error  *string `json:"error"`
}
