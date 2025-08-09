package model

import (
	"database/sql"
	"errors"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
)

func GetRatingByProdId(productId int) (*entity.Rating, error) {
	var rating entity.Rating
	err := db.DB.Get(&rating, "SELECT average_rating, rating_count FROM products WHERE id = $1", productId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &rating, nil
}

func GetRatingCommentsByProdId(productId int) ([]entity.RatingComments, error) {
	var comments []entity.RatingComments
	err := db.DB.Select(&comments, `SELECT u.name AS u_name, u.img, pv.name AS pv_name, pv.interval, i.name AS i_name,
	r.rating, r.comment, r.created_at
	FROM ratings r
	JOIN users u ON u.id = r.user_id
	JOIN product_variants pv ON pv.id = r.product_variant_id
	JOIN intervals i ON i.id = pv.interval_id
	WHERE r.product_id = $1 ORDER BY r.rating DESC, r.created_at DESC`, productId)
	if err != nil {
		return nil, err
	}

	if len(comments) == 0 {
		return nil, nil
	}

	return comments, nil
}
