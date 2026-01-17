package model

import (
	"log"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
)

func FetchReminderBatch(
	today string,
	limit int,
) ([]entity.ReminderSchedule, error) {
	var args []any
	query := `
		SELECT
			id,
			user_id,
			product_id,
			product_variant_id,
			next_send,
			next_warning_send,
			next_remove,
			last_sent_at,
			is_over
		FROM reminder_schedules
		WHERE last_sent_at IS DISTINCT FROM $1
		ORDER BY created_at
		LIMIT $2
		`
	args = []any{today, limit}

	var rows []entity.ReminderSchedule
	err := db.DB.Select(&rows, query, args...)
	log.Println(len(rows))
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
