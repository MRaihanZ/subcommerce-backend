package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/errs"
)

func GetProductsBySellerId(id interface{}) ([]*entity.JsonProductAdd, error) {
	var rows []entity.QueryProductAdd
	err := db.DB.Select(&rows, `SELECT p.id AS p_id, p.name AS p_name, p.description, p.sold, 
						p.average_rating, pi.img, p.active,
						i.id AS i_id, i.name AS i_name, 
						pv.id AS pv_id, pv.is_default, pv.name AS pv_name, pv.interval, 
						pv.stock, pv.sold AS pv_sold, pv.price, pv.discount, pv.min_order
					FROM products p
					JOIN product_images pi ON pi.product_id = p.id
					JOIN product_variants pv ON pv.product_id = p.id
					JOIN intervals i ON pv.interval_id = i.id
					JOIN sellers s ON p.seller_id = s.id
					WHERE s.id = $1 
					ORDER BY p.created_at DESC, pv_name ASC;`, id)

	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	// map to group products by product ID
	productsMap := make(map[int]*entity.JsonProductAdd)

	// map to track variants for each product
	productVariantsMap := make(map[int]map[int]bool)

	for _, row := range rows {
		// check if product already exists
		if _, exists := productsMap[row.PId]; !exists {
			productsMap[row.PId] = &entity.JsonProductAdd{
				PId:           row.PId,
				PName:         row.PName,
				Description:   row.Description,
				Sold:          row.Sold,
				AverageRating: row.AverageRating,
				Active:        row.Active,
			}
			productVariantsMap[row.PId] = make(map[int]bool)
		}

		product := productsMap[row.PId]

		// add image (avoid duplicates)
		product.PiImages = append(product.PiImages, entity.ProductImage{
			Image: row.PiImage,
		})

		// add variant if not added yet
		if !productVariantsMap[row.PId][row.PvId] {
			product.ProductVariants = append(product.ProductVariants, entity.ProductVariantAdd{
				PvId:        row.PvId,
				IId:         row.IId,
				IName:       row.IName,
				PvIsDefault: row.PvIsDefault,
				PvName:      row.PvName,
				PvInterval:  row.PvInterval,
				PvStock:     row.PvStock,
				PvSold:      row.PvSold,
				PvPrice:     row.PvPrice,
				PvDiscount:  row.PvDiscount,
				PvMinOrder:  row.PvMinOrder,
			})
			productVariantsMap[row.PId][row.PvId] = true
		}
	}

	// convert map to slice
	products := make([]*entity.JsonProductAdd, 0, len(productsMap))
	for _, product := range productsMap {
		products = append(products, product)
	}

	return products, nil
}

func GetAllProductsSummarize(search string, min int, max int) ([]entity.ProductSummarize, error) {
	baseQuery := `
		SELECT 
			p.id AS p_id,
			p.name AS p_name,
			p.sold,
			p.average_rating,
			s.name AS seller_name,
			pi.img,
			pv.id AS pv_id,
			pv.name AS pv_name,
			pv.price,
			pv.discount
		FROM products p
		JOIN sellers s ON p.seller_id = s.id
		JOIN LATERAL (
			SELECT pi.img 
			FROM product_images pi 
			WHERE pi.product_id = p.id 
			LIMIT 1
		) pi ON true
		JOIN LATERAL (
			SELECT pv.id, pv.name, pv.price, pv.discount 
			FROM product_variants pv 
			WHERE pv.product_id = p.id 
			LIMIT 1
		) pv ON true
		WHERE p.active = true
	`

	args := []interface{}{}
	argPos := 1

	if search != "" {
		baseQuery += fmt.Sprintf(" AND p.name ILIKE $%d", argPos)
		args = append(args, "%"+search+"%")
		argPos++
	}

	if min > 0 {
		baseQuery += fmt.Sprintf(" AND pv.price >= $%d", argPos)
		args = append(args, min)
		argPos++
	}

	if max > 0 {
		baseQuery += fmt.Sprintf(" AND pv.price <= $%d", argPos)
		args = append(args, max)
		argPos++
	}

	var products []entity.ProductSummarize
	err := db.DB.Select(&products, baseQuery, args...)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, nil
	}

	return products, nil
}

func GetAllProductsHotSummarize(min int, max int) ([]entity.ProductSummarize, error) {
	baseQuery := `
		SELECT 
			p.id AS p_id,
			p.name AS p_name,
			p.sold,
			p.average_rating,
			s.name AS seller_name,
			pi.img,
			pv.id AS pv_id,
			pv.name AS pv_name,
			pv.price,
			pv.discount
		FROM products p
		JOIN sellers s ON p.seller_id = s.id
		JOIN LATERAL (
			SELECT pi.img
			FROM product_images pi
			WHERE pi.product_id = p.id
			LIMIT 1
		) pi ON true
		JOIN LATERAL (
			SELECT pv.id, pv.name, pv.price, pv.discount
			FROM product_variants pv
			WHERE pv.product_id = p.id
			ORDER BY pv.sold DESC
			LIMIT 1
		) pv ON true
		WHERE p.active = true
	`

	args := []interface{}{}
	argPos := 1

	if min > 0 {
		baseQuery += fmt.Sprintf(" AND pv.price >= $%d", argPos)
		args = append(args, min)
		argPos++
	}

	if max > 0 {
		baseQuery += fmt.Sprintf(" AND pv.price <= $%d", argPos)
		args = append(args, max)
		argPos++
	}

	baseQuery += `
		ORDER BY (p.average_rating * p.sold) DESC
	`

	var products []entity.ProductSummarize
	err := db.DB.Select(&products, baseQuery, args...)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, nil
	}

	return products, nil
}

func GetAllProductsDiscountSummarize(min int, max int) ([]entity.ProductSummarize, error) {
	baseQuery := `
		SELECT 
			p.id AS p_id,
			p.name AS p_name,
			p.sold,
			p.average_rating,
			s.name AS seller_name,
			pi.img,
			pv.id AS pv_id,
			pv.name AS pv_name,
			pv.price,
			pv.discount
		FROM products p
		JOIN sellers s ON p.seller_id = s.id
		JOIN LATERAL (
			SELECT pi.img
			FROM product_images pi
			WHERE pi.product_id = p.id
			LIMIT 1
		) pi ON true
		JOIN LATERAL (
			SELECT pv.id, pv.name, pv.price, pv.discount
			FROM product_variants pv
			WHERE pv.product_id = p.id
			ORDER BY pv.discount DESC
			LIMIT 1
		) pv ON true
		WHERE p.active = true
		  AND pv.discount > 0
	`

	args := []interface{}{}
	argPos := 1

	if min > 0 {
		baseQuery += fmt.Sprintf(" AND pv.price >= $%d", argPos)
		args = append(args, min)
		argPos++
	}

	if max > 0 {
		baseQuery += fmt.Sprintf(" AND pv.price <= $%d", argPos)
		args = append(args, max)
		argPos++
	}

	baseQuery += `
		ORDER BY pv.discount DESC, p.average_rating DESC, p.sold DESC
	`

	var products []entity.ProductSummarize
	err := db.DB.Select(&products, baseQuery, args...)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, nil
	}

	return products, nil
}

func GetProductById(id string) (*entity.JsonProduct, error) {
	var products []entity.QueryProduct
	err := db.DB.Select(&products, `SELECT p.id AS p_id, p.name AS p_name, p.description, p.sold, p.average_rating, p.rating_count, p.created_at, pi.img, 
					i.code, i.name AS i_name, pv.id AS pv_id, pv.is_default, pv.name AS pv_name, pv.interval, pv.stock, pv.sold AS pv_sold,
					pv.price, pv.discount, pv.min_order FROM products p
					JOIN product_images pi ON pi.product_id = p.id
					JOIN product_variants pv ON pv.product_id = p.id
					JOIN intervals i ON pv.interval_id = i.id
					WHERE p.active = true AND p.id = $1 ORDER BY pv_name;`, id)

	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, nil
	}

	productsMap := &entity.JsonProduct{
		PId:           products[0].PId,
		PName:         products[0].PName,
		Description:   products[0].Description,
		Sold:          products[0].Sold,
		AverageRating: products[0].AverageRating,
		RatingCount:   products[0].RatingCount,
		CreatedAt:     products[0].CreatedAt,
	}

	productVariantsMap := make(map[int]*entity.ProductVariant)

	for _, row := range products {
		productsMap.PiImages = append(productsMap.PiImages, entity.ProductImage{
			Image: row.PiImage,
		})
		if _, exist := productVariantsMap[row.PvId]; !exist {
			productVariantsMap[row.PvId] = &entity.ProductVariant{
				PvId: row.PvId,
			}
			productsMap.ProductVariants = append(productsMap.ProductVariants, entity.ProductVariant{
				PvId:        row.PvId,
				ICode:       row.ICode,
				IName:       row.IName,
				PvIsDefault: row.PvIsDefault,
				PvName:      row.PvName,
				PvInterval:  row.PvInterval,
				PvStock:     row.PvStock,
				PvSold:      row.PvSold,
				PvPrice:     row.PvPrice,
				PvDiscount:  row.PvDiscount,
				PvMinOrder:  row.PvMinOrder,
			})
		}
	}

	return productsMap, nil
}

func GetProductImagePath(id interface{}, productId int) ([]string, error) {
	var imgPath []string
	err := db.DB.Select(&imgPath, `SELECT pi.img FROM product_images pi
	JOIN products p ON pi.product_id = p.id
	WHERE p.seller_id = $1 AND pi.product_id = $2;`, id, productId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return imgPath, nil
}

func DeleteProductImage(productId int) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var returnIdDelete int
	err := db.DB.QueryRowxContext(ctx, `DELETE FROM product_images WHERE product_id = $1 RETURNING product_id;`, productId).Scan(&returnIdDelete)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoProductImagesFound
		}
		return nil, err
	}
	return &returnIdDelete, nil
}

func UpdateProductImage(productId int, img []string) (*int, error) {
	query := `INSERT INTO product_images
	(product_id, img) VALUES`

	args := []interface{}{}
	reqData := []string{}

	for i, arg := range img {
		n := i*2 + 1
		reqData = append(reqData, fmt.Sprintf("($%d, $%d)", n, n+1))
		args = append(args, productId, arg)
	}
	query += strings.Join(reqData, ",")
	query += "RETURNING product_id"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnProdId int

	err := db.DB.QueryRowxContext(ctx, query, args...).Scan(&returnProdId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoProductImagesFound
		}
		return nil, err
	}

	return &returnProdId, nil
}

func UpdateProductDefault(id interface{}, productId int, productVariantId int, payload entity.ReqProductUpdate) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnId int
	err := db.DB.QueryRowxContext(ctx, `
	WITH products_update AS (
	UPDATE products
	SET name = $1, description = $2, active = $3
	WHERE seller_id = $4 AND id = $5
	RETURNING id)
	UPDATE product_variants
	SET stock = $6, discount = $7, price = $8, min_order = $9, interval = $10, interval_id = $11
	WHERE id = $12 AND product_id IN (SELECT id FROM products_update)
	RETURNING id;`,
		payload.Name, payload.Description, payload.IsActive,
		id, productId, payload.Stock, payload.Discount, payload.Price,
		payload.MinPurchase, payload.Interval, payload.IID, productVariantId).Scan(&returnId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoProductFound
		}
		return nil, err
	}
	return &returnId, nil
}

func UpdateProduct(id interface{}, productId int, payload entity.ReqProductUpdate) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnId int
	err := db.DB.QueryRowxContext(ctx, `
	UPDATE products
	SET name = $1, description = $2, active = $3
	WHERE seller_id = $4 AND id = $5
	RETURNING id;`,
		payload.Name, payload.Description, payload.IsActive, id, productId).Scan(&returnId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoProductFound
		}
		log.Println(err)
		return nil, err
	}
	return &returnId, nil
}

func GetDefaultProductVariant(id interface{}, productId int, productVariantId int) (*string, error) {
	var name string
	err := db.DB.Get(&name, `SELECT pv.name FROM product_variants pv
	JOIN products p ON pv.product_id = p.id
	WHERE p.id = $1 AND p.seller_id = $2 AND pv.id = $3 LIMIT 1`, productId, id, productVariantId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoProductVariantFound
		}
		return nil, err
	}

	return &name, nil
}

func CreateProduct(id interface{}, payload entity.ReqProductAdd) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnProdId int

	err := db.DB.QueryRowxContext(ctx, `INSERT INTO products
	(seller_id, name, description, active) VALUES
	($1, $2, $3, $4) RETURNING id`, id, payload.Name, payload.Description, payload.Active).Scan(&returnProdId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoProductFound
		}
		return nil, err
	}
	return &returnProdId, err
}

func CreateProductVariantsAdd(id interface{}, productId int, payload []entity.ReqProductVariants) (*int, error) {
	if payload[0].IsDefault {
		query := `INSERT INTO product_variants 
	(product_id, interval_id, is_default, name, interval, stock, price, discount, min_order) VALUES`

		args := []interface{}{}
		reqData := []string{}

		for i, arg := range payload {
			n := i*9 + 1
			reqData = append(reqData, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", n, n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8))
			args = append(args, productId, arg.IId, true, "default", arg.Interval, arg.Stock, arg.Price, arg.Discount, arg.MinOrder)
		}
		query += strings.Join(reqData, ",")
		query += "RETURNING product_id"

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var returnProdId int

		err := db.DB.QueryRowxContext(ctx, query, args...).Scan(&returnProdId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errs.ErrNoProductVariantFound
			}
			return nil, err
		}
		return &returnProdId, nil
	} else {
		query := `INSERT INTO product_variants 
	(product_id, interval_id, is_default, name, interval, stock, price, discount, min_order) VALUES`

		args := []interface{}{}
		reqData := []string{}

		for i, arg := range payload {
			n := i*9 + 1
			reqData = append(reqData, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", n, n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8))
			args = append(args, productId, arg.IId, false, arg.Name, arg.Interval, arg.Stock, arg.Price, arg.Discount, arg.MinOrder)
		}
		query += strings.Join(reqData, ",")
		query += "RETURNING product_id"

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var returnProdId int

		err := db.DB.QueryRowxContext(ctx, query, args...).Scan(&returnProdId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errs.ErrNoProductVariantFound
			}
			return nil, err
		}
		return &returnProdId, nil
	}
}

func CreateProductImages(id interface{}, productId int, img []string) (*int, error) {
	query := `INSERT INTO product_images 
	(product_id, img) VALUES`

	args := []interface{}{}
	reqData := []string{}

	for i, arg := range img {
		n := i*2 + 1
		reqData = append(reqData, fmt.Sprintf("($%d, $%d)", n, n+1))
		args = append(args, productId, arg)
	}
	query += strings.Join(reqData, ",")
	query += "RETURNING product_id"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnProdId int

	err := db.DB.QueryRowxContext(ctx, query, args...).Scan(&returnProdId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoProductImagesFound
		}
		return nil, err
	}
	return &returnProdId, nil
}

func CreateProductVariants(id interface{}, productId int, productVariantId int, payload []entity.ReqProductVariantUpdate) (*int, error) {
	name, err := GetDefaultProductVariant(id, productId, productVariantId)
	if err != nil {
		if errors.Is(err, errs.ErrNoProductVariantFound) {
			return nil, errs.ErrNoProductVariantFound
		}
		return nil, err
	}

	query := `INSERT INTO product_variants 
	(product_id, interval_id, is_default, name, interval, stock, price, discount, min_order) VALUES`

	args := []interface{}{}
	reqData := []string{}

	for i, arg := range payload {
		n := i*9 + 1
		reqData = append(reqData, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", n, n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8))
		args = append(args, productId, arg.IId, false, arg.Name, arg.Interval, arg.Stock, arg.Price, arg.Discount, arg.MinOrder)
	}
	query += strings.Join(reqData, ",")
	query += "RETURNING product_id"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnProdId int

	err = db.DB.QueryRowxContext(ctx, query, args...).Scan(&returnProdId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoProductVariantFound
		}
		return nil, err
	}

	var returnIdDelete int
	if *name == "default" {
		err := db.DB.QueryRowxContext(ctx, `DELETE FROM product_variants WHERE id = $1 RETURNING id;`, productVariantId).Scan(&returnIdDelete)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errs.ErrNoProductVariantFound
			}
			return nil, err
		}
	}
	return &returnProdId, nil
}

func UpdateProductVariant(id interface{}, productId int, productVariantId int, payload entity.ReqProductVariantUpdate) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnId int
	err := db.DB.QueryRowxContext(ctx, `
	UPDATE product_variants pv SET 
	interval_id = $1, name = $2, interval = $3, stock = $4,
	price = $5, discount = $6, min_order = $7
	FROM products p
	WHERE p.seller_id = $8 AND pv.id = $9 AND pv.product_id = $10 RETURNING pv.id;`,
		payload.IId, payload.Name, payload.Interval, payload.Stock, payload.Price, payload.Discount, payload.MinOrder,
		id, productVariantId, productId).Scan(&returnId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoProductFound
		}
		return nil, err
	}
	return &returnId, nil
}

func GetCountProductVariant(productId int) (*int, error) {
	var num int
	err := db.DB.Get(&num, `SELECT COUNT(*)
	FROM product_variants
	WHERE product_id = $1;`, productId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoProductVariantFound
		}
		return nil, err
	}

	return &num, nil
}

func DeleteProduct(productId int) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var returnIdDelete int
	err := db.DB.QueryRowxContext(ctx, `DELETE FROM products WHERE id = $1 RETURNING id;`, productId).Scan(&returnIdDelete)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoProductFound
		}
		return nil, err
	}
	return &returnIdDelete, nil
}

func DeleteProductVariant(productId int, productVariantId int) (*int, error) {
	num, err := GetCountProductVariant(productId)
	if err != nil {
		if errors.Is(err, errs.ErrNoProductVariantFound) {
			return nil, errs.ErrNoProductVariantFound
		}
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnIdDelete int
	if *num > 1 {
		err := db.DB.QueryRowxContext(ctx, `DELETE FROM product_variants
		WHERE id = $1 AND product_id = $2 RETURNING id;`, productVariantId, productId).Scan(&returnIdDelete)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errs.ErrNoProductVariantFound
			}
			return nil, err
		}
	} else {
		err := db.DB.QueryRowxContext(ctx, `UPDATE product_variants SET is_default = true, name = 'default'
		WHERE id = $1 AND product_id = $2 RETURNING id;`, productVariantId, productId).Scan(&returnIdDelete)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errs.ErrNoProductVariantFound
			}
			return nil, err
		}
	}
	return &returnIdDelete, nil
}
