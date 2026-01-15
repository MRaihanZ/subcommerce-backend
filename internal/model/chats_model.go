package model

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/google/uuid"
)

func GetAllConversations(queryId interface{}, stateQuery string) ([]entity.Conversation, error) {
	var rows []entity.Conversation
	var query string

	switch stateQuery {
	case "user":
		query = `SELECT c.id, c.user_id, c.seller_id, c.last_message_at, c.last_message_content, s.name, s.img
		FROM conversations c
		JOIN sellers s ON c.seller_id = s.id
		WHERE c.user_id = $1 ORDER BY last_message_at DESC`
	case "seller":
		query = `SELECT c.id, c.user_id, c.seller_id, c.last_message_at, c.last_message_content, u.name, u.img
		FROM conversations c
		JOIN users u ON c.user_id = u.id
		WHERE c.seller_id = $1 ORDER BY last_message_at DESC`
	default:
		return nil, errors.New("invalid state query")
	}
	err := db.DB.Select(&rows, query, queryId)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return rows, nil
}

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

func GetMessagesByConvID(convID string,
	before string,
	beforeID string,
	limit int) ([]entity.Message, error) {
	var rows []entity.Message
	var query string
	var args []any
	if before == "" {
		// Initial load (latest messages)
		query = `
			SELECT id, is_user, content, sent_at
			FROM messages
			WHERE conversation_id = $1
			ORDER BY sent_at DESC, id DESC
			LIMIT $2;
		`
		args = []any{convID, limit}
	} else {
		// Load older messages
		query = `
			SELECT id, is_user, content, sent_at
			FROM messages
			WHERE conversation_id = $1
			  AND (sent_at, id) < ($2, $3)
			ORDER BY sent_at DESC, id DESC
			LIMIT $4;
		`
		args = []any{convID, before, beforeID, limit}
	}
	if err := db.DB.Select(&rows, query, args...); err != nil {
		return nil, err
	}

	// reverse so UI shows oldest → newest
	slices.Reverse(rows)

	return rows, nil
}

func CreateMessage(convId, senderId interface{}, isUser bool, content string) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var convIdReturn string

	err := db.DB.QueryRowxContext(ctx, `INSERT INTO messages (conversation_id, sender_id, is_user, content)
		VALUES ($1, $2, $3, $4) RETURNING conversation_id`,
		convId, senderId, isUser, content,
	).Scan(&convIdReturn)
	if err != nil {
		return nil, err
	}
	return &convIdReturn, nil
}
