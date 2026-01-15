package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/google/uuid"
)

func GetConversationsById(userId, sellerId interface{}) (*string, error) {
	var convId string

	err := db.DB.Get(&convId, `SELECT id FROM conversations WHERE user_id = $1 AND seller_id = $2`, userId, sellerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &convId, nil
}

func CreateConversation(userId, sellerId interface{}) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id uuid.UUID
	var convId string

	id = uuid.New()
	err := db.DB.QueryRowxContext(ctx, "INSERT INTO conversations (id, user_id, seller_id) VALUES ($1, $2, $3) RETURNING id",
		id, userId, sellerId,
	).Scan(&convId)
	if err != nil {
		return nil, err
	}
	return &convId, nil
}
