package model

import (
	"context"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/google/uuid"
)

func CreatePayout(sellerId interface{}, name, payId, description string, amount int64) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id uuid.UUID
	var returnId string
	id = uuid.New()

	query := `INSERT INTO payout (id, seller_id, transfer_name, transfer_type, transfer_id, transfer_amount, transfer_description, transfer_status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	args := []interface{}{
		id,
		sellerId,
		name,
		"gopay",
		payId,
		amount,
		description,
		"permintaan sedang dalam antrian",
	}

	err := db.DB.QueryRowxContext(ctx, query, args...).Scan(&returnId)
	if err != nil {
		return nil, err
	}
	return &returnId, nil
}

func GetPayouts(
	ctx context.Context,
) ([]entity.Payout, error) {

	var payouts []entity.Payout

	query := `
		SELECT
			p.id,
			s.name AS seller_name,
			p.created_at,
			p.transfer_name,
			p.transfer_type,
			p.transfer_id,
			p.transfer_amount,
			p.transfer_description,
			p.transfer_status
		FROM payout p
		JOIN sellers s ON s.id = p.seller_id
		ORDER BY p.created_at DESC, p.transfer_status ASC
	`

	err := db.DB.SelectContext(ctx, &payouts, query)
	return payouts, err
}

func UpdatePayoutStatus(
	ctx context.Context,
	payoutID string,
	status string,
) error {

	query := `
		UPDATE payout
		SET transfer_status = $1
		WHERE id = $2
	`

	_, err := db.DB.ExecContext(ctx, query, status, payoutID)
	return err
}
