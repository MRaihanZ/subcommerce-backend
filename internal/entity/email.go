package entity

type Sender interface {
	SendMail(to, subject, html string) error
}

type SubscriptionEmailData struct {
	OrderID         string
	SubscriptionID  string
	UserName        string
	ProductName1    string
	ProductName2    string
	Subscription    string
	Status          string
	OrderDate       string
	PaymentDeadline string
	TargetEmail     string
	Domain          string
}

type CancelationSubscriptionByUserEmailData struct {
	SellerName         string `json:"seller_name"`
	UserName           string
	ProductName        string `json:"product_name"`
	ProductVariantName string `json:"product_variant_name"`
	Timestamp          string
	Domain             string
}
