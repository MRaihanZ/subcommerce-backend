package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/MRaihanZ/subcommerce-backend/internal/errs"
	"golang.org/x/crypto/bcrypt"
)

func GetUserById(id interface{}) (*entity.User, error) {
	var user entity.User
	err := db.DB.Get(&user, "SELECT name, img, email, dob, created_at FROM users WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func UpdateUser(id interface{}, name, img, email, password string, dob time.Time) (*entity.User, error) {
	var user entity.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		err = db.DB.QueryRowxContext(ctx, "UPDATE users SET name = $1, img = $2, email = $3, dob = $4, password = $5 WHERE id = $6 RETURNING name, img, email, dob, created_at",
			name, img, email, dob, hashedPassword, id).Scan(&user.Name, &user.Img, &user.Email, &user.Dob, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
	} else {
		err := db.DB.QueryRowxContext(ctx, "UPDATE users SET name = $1, img = $2, email = $3, dob = $4 WHERE id = $5 RETURNING name, img, email, dob, created_at",
			name, img, email, dob, id).Scan(&user.Name, &user.Img, &user.Email, &user.Dob, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func GetUserImagePath(id interface{}) (*string, error) {
	var imgPath string
	err := db.DB.Get(&imgPath, "SELECT img FROM users WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &imgPath, nil
}

func DeleteUser(id interface{}) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var deletedName string
	err := db.DB.QueryRowxContext(ctx, `DELETE FROM users WHERE id = $1
	RETURNING name`, id).Scan(&deletedName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return &deletedName, nil
}
