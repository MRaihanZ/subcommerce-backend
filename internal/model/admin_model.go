package model

import (
	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
)

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
