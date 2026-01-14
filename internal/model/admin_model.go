package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/errs"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetAllAdmin() ([]entity.Admin, error) {
	var rows []entity.Admin
	err := db.DB.Select(&rows, "SELECT id, name, email FROM admins")
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return rows, nil
}

func GetAdminByName(name string) ([]entity.Admin, error) {
	var rows []entity.Admin
	err := db.DB.Select(&rows, `SELECT id, name, email FROM admins WHERE name ILIKE $1`, "%"+name+"%")
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return rows, nil
}

func CreateAdmin(req entity.CreateAdmin) (*string, error) {
	checkEmail, err := GetAdminByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if checkEmail != nil {
		return nil, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id uuid.UUID
	var returnId string
	id = uuid.New()
	err = db.DB.QueryRowxContext(ctx, "INSERT INTO admins (id, name, email, password) VALUES ($1, $2, $3, $4) RETURNING id",
		id, req.Name, req.Email, hashedPassword,
	).Scan(&returnId)
	if err != nil {
		return nil, err
	}
	return &returnId, nil
}

func UpdateAdmin(id interface{}, name, email, password string) (*entity.Admin, error) {
	var admin entity.Admin

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		err = db.DB.QueryRowxContext(ctx, "UPDATE admins SET name = $1, email = $2, password = $3 WHERE id = $4 RETURNING id, name, email",
			name, email, hashedPassword, id).Scan(&admin.Id, &admin.Name, &admin.Email)
		if err != nil {
			return nil, err
		}
	} else {
		err := db.DB.QueryRowxContext(ctx, "UPDATE admins SET name = $1, email = $2 WHERE id = $3 RETURNING id, name, email",
			name, email, id).Scan(&admin.Id, &admin.Name, &admin.Email)
		if err != nil {
			return nil, err
		}
	}

	return &admin, nil
}

func DeleteAdmin(id interface{}) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var deletedName string
	err := db.DB.QueryRowxContext(ctx, `DELETE FROM admins WHERE id = $1
	RETURNING name`, id).Scan(&deletedName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return &deletedName, nil
}

func GetAllUsers() ([]entity.Users, error) {
	var rows []entity.Users
	err := db.DB.Select(&rows, "SELECT id, name, img, email, dob, created_at FROM users")
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return rows, nil
}

func GetUsersByName(name string) ([]entity.Users, error) {
	var rows []entity.Users
	err := db.DB.Select(&rows, `SELECT id, name, img, email, dob, created_at FROM users WHERE name ILIKE $1`, "%"+name+"%")
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return rows, nil
}

func GetAllSellers() ([]entity.Seller, error) {
	var rows []entity.Seller
	err := db.DB.Select(&rows, `SELECT id, name, img, address, created_at FROM sellers;`)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return rows, nil
}

func GetSellersByName(name string) ([]entity.Seller, error) {
	var rows []entity.Seller
	err := db.DB.Select(&rows, `SELECT id, name, img, address, created_at FROM sellers WHERE name ILIKE $1`, "%"+name+"%")
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return rows, nil
}

func GetAllProducts() ([]*entity.JsonProductAdd, error) {
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
					ORDER BY p.name ASC;`)

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

func GetProductsByName(name string) ([]*entity.JsonProductAdd, error) {
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
					WHERE p.name ILIKE $1
					ORDER BY p.name ASC;`, "%"+name+"%")

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

func UpdateActiveProduct(productId int, active bool) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnId int
	err := db.DB.QueryRowxContext(ctx, `
	UPDATE products
	SET active = $1
	WHERE id = $2
	RETURNING id`, active, productId).Scan(&returnId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNoProductFound
		}
		return nil, err
	}
	return &returnId, nil
}
