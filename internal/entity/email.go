package entity

type Sender interface {
	SendMail(to, subject, html string) error
}

type SubscriptionEmailData struct {
	ID              string
	OrderID         int
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
