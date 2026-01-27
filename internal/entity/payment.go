package entity

import "time"

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

type CreatePayoutRequest struct {
	Name        string `json:"name"`
	PayId       string `json:"pay_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

type Payout struct {
	ID                  string    `db:"id" json:"id"`
	SellerName          string    `db:"seller_name" json:"seller_name"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	TransferName        string    `db:"transfer_name" json:"transfer_name"`
	TransferType        string    `db:"transfer_type" json:"transfer_type"`
	TransferID          string    `db:"transfer_id" json:"transfer_id"`
	TransferAmount      int64     `db:"transfer_amount" json:"transfer_amount"`
	TransferDescription string    `db:"transfer_description" json:"transfer_description"`
	TransferStatus      *string   `db:"transfer_status" json:"transfer_status"`
}

type UpdatePayoutStatusRequest struct {
	Status string `json:"status"`
}
