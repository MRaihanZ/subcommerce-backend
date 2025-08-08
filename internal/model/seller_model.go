package model

import (
	"database/sql"
	"errors"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
)

func GetSellerById(id string) (*entity.SellerSummarize, error) {
	var seller entity.SellerSummarize
	err := db.DB.Get(&seller, "SELECT id, img, sold_products, average_rating FROM sellers WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &seller, nil
}
