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
