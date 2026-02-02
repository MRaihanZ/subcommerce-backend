package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/errs"
)

func CreateRating(productId int, productVariantId int, orderId int, userId interface{}, req *entity.RatingRequest) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var ratingCheck bool
	err := db.DB.GetContext(ctx, &ratingCheck, `SELECT rating FROM orders
	WHERE id = $1 AND product_id = $2 AND product_variant_id = $3;
	`, orderId, productId, productVariantId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoRatingFound
		}
		return nil, err
	}

	if ratingCheck {
		return nil, errs.ErrAlreadyRated
	}

	var rating int

	err = db.DB.QueryRowContext(ctx, "INSERT INTO ratings (product_id, product_variant_id, user_id, order_id, rating, comment) VALUES ($1, $2, $3, $4, $5, $6) RETURNING rating",
		productId, productVariantId, userId, orderId, req.Rating, req.Comment).Scan(&rating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoRatingFound
		}
		return nil, err
	}

	_, err = UpdateRatingOrder(productId, productVariantId, orderId, userId)
	if err != nil {
		if errors.Is(err, errs.ErrNoOrderFound) {
			return nil, errs.ErrNoOrderFound
		}
		return nil, err
	}
	return &rating, nil
}

func UpdateRatingProductSeller(productId int, req *entity.RatingRequest) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnId string
	err := db.DB.QueryRowContext(ctx, `
	WITH updated_product AS (
	UPDATE products SET rating_count = rating_count + 1, rating_total = rating_total + $1, average_rating = ((rating_total + $1)::numeric / (rating_count + 1))::numeric(4,3)
	WHERE id = $2
	RETURNING seller_id
	)
	
	UPDATE sellers SET rating_count = rating_count + 1, rating_total = rating_total + $1, average_rating = ((rating_total + $1)::numeric / (rating_count + 1))::numeric(4,3)
	WHERE id = (
		SELECT seller_id
		FROM updated_product
	) RETURNING id;
	`, req.Rating, productId).Scan(&returnId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoOrderFound
		}
		return nil, err
	}

	return &returnId, nil
}

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
