package model

import (
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
)

func FetchReminderBatch(
	today string,
	limit int,
) ([]entity.SubscriptionData, error) {
	var args []any
	query := `SELECT
		rs.id, rs.user_id, rs.product_id, rs.product_variant_id, rs.next_send, rs.next_warning_send, rs.next_remove, rs.last_sent_at, rs.is_over,
		o.order_pretty_id, o.created_at, u.name AS user_name, u.email, p.name AS product_name, pv.name AS product_variant_name, pv.interval, i.name AS interval_name
	FROM reminder_schedules rs
	JOIN orders o ON o.order_uq_id = rs.id
	JOIN users u ON u.id = rs.user_id
	JOIN products p ON p.id = rs.product_id
	JOIN product_variants pv ON pv.id = rs.product_variant_id
	JOIN intervals i ON i.id = pv.interval_id
	WHERE last_sent_at IS DISTINCT FROM $1
	AND (
		next_send::date = $1
		OR next_warning_send::date = $1
		OR next_remove::date = $1
	)
	ORDER BY created_at
	LIMIT $2;`
	args = []any{today, limit}

	var rows []entity.SubscriptionData
	err := db.DB.Select(&rows, query, args...)
	return rows, err
}

func DeleteReminder(id string) error {
	_, err := db.DB.Exec(
		`DELETE FROM reminder_schedules WHERE id = $1`,
		id,
	)
	return err
}

func UpdateLastSentAt(id string, t time.Time) error {
	_, err := db.DB.Exec(
		`UPDATE reminder_schedules SET last_sent_at = $1 WHERE id = $2`,
		t,
		id,
	)
	return err
}

func UpdateIsOver(id string) error {
	_, err := db.DB.Exec(
		`UPDATE reminder_schedules SET is_over = true WHERE id = $1`,
		id,
	)
	return err
}
