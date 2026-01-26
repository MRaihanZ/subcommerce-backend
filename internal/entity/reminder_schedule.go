package entity

import "time"

type SubscriptionData struct {
	ID               string    `db:"id"`
	UserID           string    `db:"user_id"`
	ProductID        int       `db:"product_id"`
	ProductVariantID int       `db:"product_variant_id"`
	NextSend         time.Time `db:"next_send"`
	NextWarningSend  time.Time `db:"next_warning_send"`
	NextRemove       time.Time `db:"next_remove"`
	LastSentAt       time.Time `db:"last_sent_at"`
	IsOver           bool      `db:"is_over"`
	OrderPrettyID    string    `db:"order_pretty_id"`
	CreatedAt        time.Time `db:"created_at"`
	UName            string    `db:"user_name"`
	UEmail           string    `db:"email"`
	PName            string    `db:"product_name"`
	PVName           string    `db:"product_variant_name"`
	PVInterval       int       `db:"interval"`
	IName            string    `db:"interval_name"`
}
