package model

import (
	"database/sql"
	"errors"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(name string, email string, dob string, password string) (*string, error) {
	checkEmail, err := GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if checkEmail != nil {
		return nil, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var id uuid.UUID
	var returnId string
	id = uuid.New()
	err = db.DB.QueryRowx("INSERT INTO users (id, name, email, dob, password) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		id, name, email, dob, hashedPassword,
	).Scan(&returnId)
	if err != nil {
		return nil, err
	}
	return &returnId, nil
}

func GetUserByEmail(email string) (*entity.SignIn, error) {
	var signIn entity.SignIn
	err := db.DB.Get(&signIn, "SELECT id, email, password FROM users WHERE email = $1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &signIn, nil
}

func GetSellerByEmail(id string) (*string, error) {
	var signIn string
	err := db.DB.Get(&signIn, "SELECT id FROM sellers WHERE user_id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &signIn, nil
}
