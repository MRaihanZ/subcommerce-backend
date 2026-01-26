package model

import (
	"context"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/google/uuid"
)

func CreatePayout(sellerId interface{}, name, payId, description string, amount int64) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id uuid.UUID
	var returnId string
	id = uuid.New()

	query := `INSERT INTO payout (id, seller_id, transfer_name, transfer_type, transfer_id, transfer_amount, transfer_description, transfer_status, payout_pretty_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'POUT-' || TO_CHAR(NOW(), 'YYYYMMDD') || '-' ||
		LPAD(nextval('orders_order_pretty_id_seq')::text, 4, '0')) RETURNING id`
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
