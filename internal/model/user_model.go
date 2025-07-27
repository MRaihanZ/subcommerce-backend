package model

import (
	"database/sql"
	"errors"
	"log"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers() ([]entity.User, error) {
	var check []entity.User
	err := db.DB.Select(&check, "SELECT * FROM users")
	if err != nil {
		log.Println("ERROR")
		return nil, err
	}

	if len(check) == 0 {
		log.Println("NO DATA")
		return nil, nil
	}

	return check, nil
}

func GetUserById(id string) (*entity.User, error) {
	var user entity.User
	err := db.DB.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func CreateUser(name string, email string, dob string, password string) (*uuid.UUID, error) {
	var id uuid.UUID
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id = uuid.New()
	err = db.DB.QueryRowx("INSERT INTO users (id, name, email, dob, password) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		id, name, email, dob, hashedPassword,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
