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

func GetAllCartProduct(id interface{}) ([]entity.CartProduct, error) {
	var carts []entity.CartProduct
	err := db.DB.Select(&carts, `SELECT s.name AS s_name,p.id AS p_id, p.name AS p_name, pi.img, p.active,
						pv.id AS pv_id, pv.name AS pv_name, pv.interval, pv.stock,
						i.name AS i_name, c.quantity, pv.price, pv.min_order
						FROM carts c
						JOIN products p ON c.product_id = p.id
						JOIN sellers s ON p.seller_id = s.id
						JOIN product_variants pv ON c.product_variant_id = pv.id
						JOIN intervals i ON pv.interval_id = i.id
						JOIN LATERAL (SELECT pi.img FROM product_images pi WHERE pi.product_id = p.id LIMIT 1) pi ON true
						WHERE c.user_id = $1 ORDER BY p.active DESC, p.created_at DESC;`, id)
	if err != nil {
		return nil, err
	}

	if len(carts) == 0 {
		return nil, nil
	}

	return carts, nil
}

func CreateCart(id interface{}, productId int, productVariantId int, quantity int) (*int, error) {
	if quantity <= 0 {
		return nil, errs.ErrQuantityLessThanZero
	}

	var min_order int
	err := db.DB.Get(&min_order, "SELECT min_order FROM product_variants WHERE id = $1 AND product_id = $2", productVariantId, productId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrProductNotFound
		}
		return nil, err
	}

	if quantity < min_order {
		return nil, errs.ErrQuantityLessThanMinOrder
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnId int
	err = db.DB.QueryRowxContext(ctx, `INSERT INTO carts(user_id, product_id, product_variant_id, quantity)
	SELECT $1, p.id, pv.id, $4
	FROM products p
	JOIN product_variants pv ON pv.product_id = p.id
	WHERE p.id = $2 AND pv.id = $3 AND p.active = true AND pv.stock >= $4
	ON CONFLICT (user_id, product_id, product_variant_id) DO UPDATE
	SET quantity = carts.quantity + EXCLUDED.quantity
	WHERE carts.quantity + EXCLUDED.quantity <= (
	SELECT stock FROM product_variants WHERE id = EXCLUDED.product_variant_id)
	RETURNING id;`, id, productId, productVariantId, quantity).Scan(&returnId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotEnoughStock
		}
		return nil, err
	}
	return &returnId, nil
}

func UpdateCart(id interface{}, productId int, productVariantId int, quantity int) (*int, error) {
	if quantity <= 0 {
		return nil, errs.ErrQuantityLessThanZero
	}

	var min_order int
	err := db.DB.Get(&min_order, "SELECT min_order FROM product_variants WHERE id = $1 AND product_id = $2", productVariantId, productId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrProductNotFound
		}
		return nil, err
	}

	if quantity < min_order {
		return nil, errs.ErrQuantityLessThanMinOrder
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnId int
	err = db.DB.QueryRowxContext(ctx, `UPDATE carts c SET quantity = $4
	FROM product_variants pv
	WHERE c.user_id = $1 AND c.product_id = $2 AND c.product_variant_id = $3
	AND pv.id = c.product_variant_id AND $4 <= pv.stock
	RETURNING c.id`, id, productId, productVariantId, quantity).Scan(&returnId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotEnoughStock
		}
		return nil, err
	}
	return &returnId, nil
}

func DeleteCarts(id interface{}) (*int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.DB.ExecContext(ctx, "DELETE FROM carts WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, errs.ErrCartsEmpty
	}

	return &rows, nil
}

func DeleteCart(id interface{}, prodId int, prodVarId int) (*entity.DeleteProductCart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var DeleteProduct entity.DeleteProductCart
	err := db.DB.QueryRowxContext(ctx, `DELETE FROM carts
	WHERE user_id = $1 AND product_id = $2 AND product_variant_id = $3
	RETURNING product_id, product_variant_id`, id, prodId, prodVarId).StructScan(&DeleteProduct)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoCartProduct
		}
		return nil, err
	}

	return &DeleteProduct, nil
}
