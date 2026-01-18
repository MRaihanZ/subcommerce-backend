package model

import (
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
)

type Payment struct {
	ID         string
	OrderID    string
	Amount     int64
	Status     string
	PaymentURL string
	ExpiredAt  time.Time
	CreatedAt  time.Time
}

func CreatePayment(p *Payment) error {
	// example using sqlx
	query := `
		INSERT INTO payments 
		(id, order_id, amount, status, payment_url, expired_at, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`
	_, err := db.DB.Exec(
		query,
		p.ID,
		p.OrderID,
		p.Amount,
		p.Status,
		p.PaymentURL,
		p.ExpiredAt,
		p.CreatedAt,
	)
	return err
}

func UpdatePaymentStatus(orderID string, status string) error {
	query := `
		UPDATE payments 
		SET status = $1 
		WHERE order_id = $2
	`
	_, err := db.DB.Exec(query, status, orderID)
	return err
}
