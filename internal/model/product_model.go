package model

import (
	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
)

func GetAllProductsSummarize() ([]entity.ProductSummarize, error) {
	var products []entity.ProductSummarize
	err := db.DB.Select(&products, `SELECT p.id AS p_id, p.name AS p_name, p.sold, p.average_rating, s.name AS seller_name,
						pi.img, pv.id AS pv_id, pv.name AS pv_name, pv.price, pv.discount
						FROM products p 
						JOIN sellers s ON p.seller_id = s.id
						JOIN LATERAL (SELECT pi.img FROM product_images pi WHERE pi.product_id = p.id LIMIT 1) pi ON true
						JOIN LATERAL (SELECT pv.id, pv.name, pv.price, pv.discount FROM product_variants pv WHERE pv.product_id = p.id LIMIT 1) pv ON true
						WHERE p.active = true;`)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, nil
	}

	return products, nil
}

func GetAllProductsHotSummarize() ([]entity.ProductSummarize, error) {
	var products []entity.ProductSummarize
	err := db.DB.Select(&products, `SELECT p.id AS p_id, p.name AS p_name, p.sold, p.average_rating, 
						s.name AS seller_name, pi.img, pv.id AS pv_id, pv.name AS pv_name, pv.price, pv.discount
						FROM products p 
						JOIN sellers s ON p.seller_id = s.id
						JOIN LATERAL (SELECT pi.img FROM product_images pi WHERE pi.product_id = p.id LIMIT 1) pi ON true
						JOIN LATERAL (SELECT pv.id, pv.name, pv.price, pv.discount FROM product_variants pv 
						WHERE pv.product_id = p.id ORDER BY pv.sold DESC LIMIT 1) pv ON true
						WHERE p.active = true ORDER BY (p.average_rating * p.sold) DESC;`)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, nil
	}

	return products, nil
}

func GetAllProductsDiscountSummarize() ([]entity.ProductSummarize, error) {
	var products []entity.ProductSummarize
	err := db.DB.Select(&products, `SELECT p.id AS p_id, p.name AS p_name, p.sold, p.average_rating, s.name AS seller_name,
						pi.img, pv.id AS pv_id, pv.name AS pv_name, pv.price, pv.discount 
						FROM products p 
						JOIN sellers s ON p.seller_id = s.id
						JOIN LATERAL (SELECT pi.img FROM product_images pi WHERE pi.product_id = p.id LIMIT 1) pi ON true
						JOIN LATERAL (SELECT pv.id, pv.name, pv.price, pv.discount FROM product_variants pv WHERE pv.product_id = p.id ORDER BY pv.discount DESC LIMIT 1) pv ON true
						WHERE p.active = true AND pv.discount > 0 
						ORDER BY pv.discount DESC, p.average_rating DESC, p.sold DESC;`)
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
