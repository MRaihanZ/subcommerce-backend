package entity

import "time"

type ReminderSchedule struct {
	ID               string    `db:"id"`
	UserID           string    `db:"user_id"`
	ProductID        int       `db:"product_id"`
	ProductVariantID int       `db:"product_variant_id"`
	NextSend         time.Time `db:"next_send"`
	NextWarningSend  time.Time `db:"next_warning_send"`
	NextRemove       time.Time `db:"next_remove"`
	LastSentAt       time.Time `db:"last_sent_at"`
	IsOver           bool      `db:"is_over"`
}
