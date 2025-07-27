package model

import (
	"database/sql"
	"errors"
	"log"

	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
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
