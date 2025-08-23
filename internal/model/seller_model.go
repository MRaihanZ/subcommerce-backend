package model

import (
	"database/sql"
	"errors"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
)

func GetSellerById(id interface{}) (*entity.Seller, error) {
	var seller entity.Seller
	err := db.DB.Get(&seller, `SELECT id, name, img, address, total_sold_products,
	average_rating, rating_total, rating_count, current_month_sales, previous_month_sales,
	current_month_cancellations, previous_month_cancellations, current_month_revenue,
	previous_month_revenue, current_month, previous_month, created_at
	FROM sellers WHERE id =  $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &seller, nil
}

func GetSellerByIdSummarize(id int) (*entity.SellerSummarize, error) {
	var seller entity.SellerSummarize
	err := db.DB.Get(&seller, `SELECT s.id, s.name, s.img, s.sold_products, s.average_rating 
	FROM sellers s
	JOIN products p ON p.seller_id = s.id
	WHERE p.id = $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &seller, nil
}
