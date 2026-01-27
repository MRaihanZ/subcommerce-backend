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

func GetSellerById(id interface{}) (*entity.Seller, error) {
	var seller entity.Seller
	err := db.DB.Get(&seller, `SELECT id, name, img, address, wallet, total_sold_products,
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
	err := db.DB.Get(&seller, `SELECT s.id, s.name, s.img, s.total_sold_products, s.average_rating 
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

func UpdateSeller(id interface{}, name, img, address string) (*entity.Seller, error) {
	var seller entity.Seller

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.DB.QueryRowxContext(ctx, "UPDATE sellers SET name = $1, img = $2, address = $3 WHERE id = $4 RETURNING name, img, address",
		name, img, address, id).Scan(&seller.Name, &seller.Img, &seller.Address)
	if err != nil {
		return nil, err
	}

	return &seller, nil
}

func GetSellerImagePath(id interface{}) (*string, error) {
	var imgPath string
	err := db.DB.Get(&imgPath, "SELECT img FROM sellers WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &imgPath, nil
}

func DeleteSeller(id interface{}) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var deletedName string
	err := db.DB.QueryRowxContext(ctx, `DELETE FROM sellers WHERE id = $1
	RETURNING name`, id).Scan(&deletedName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrSellerNotFound
		}
		return nil, err
	}

	return &deletedName, nil
}

func GetCurrentWalletAmmount(sellerId interface{}) (*int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wallet int64
	err := db.DB.GetContext(ctx, &wallet, "SELECT wallet FROM sellers WHERE id = $1", sellerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &wallet, nil
}

func UpdateWalletSellerById(sellerId interface{}, amount int64) (*string, error) {
	var seller string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.DB.QueryRowxContext(ctx, "UPDATE sellers SET wallet = $1 WHERE id = $2 RETURNING id",
		amount, sellerId).Scan(&seller)
	if err != nil {
		return nil, err
	}

	return &seller, nil
}

func GetSellerNameEmailById(sellerId interface{}) (*entity.SellerNameEmail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var data entity.SellerNameEmail
	err := db.DB.GetContext(ctx, &data, `SELECT s.name, u.email
	FROM sellers s
	JOIN users u ON u.id = s.user_id
	WHERE s.id = $1;`, sellerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}
