package entity

type CreatePaymentRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	Amount  int64  `json:"amount" binding:"required"`
}

type MidtransChargeRequest struct {
	EnablePayments     []string `json:"enabled_payments"`
	TransactionDetails struct {
		OrderID  string `json:"order_id"`
		GrossAmt int64  `json:"gross_amount"`
	} `json:"transaction_details"`
	Expiry struct {
		Unit     string `json:"unit"`
		Duration int    `json:"duration"`
	} `json:"expiry"`
	CustomerRequired bool `json:"customer_required"`
	CustomerDetails  struct {
		FirstName string `json:"first_name"`
		Email     string `json:"email"`
	} `json:"customer_details"`
	ItemDetails []struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
		Price    int    `json:"price"`
	} `json:"item_details"`
	QrisDetail struct {
		Acquirer string `json:"acquirer"`
	} `json:"qris"`
	Callbacks struct {
		Finish string `json:"finish"`
	} `json:"callbacks"`
}
