package model

import (
	"github.com/MRaihanZ/subcommerce-backend/internal/db"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Test int    `json:"test"`
}

// func errorGetAllUsers() {
// 	log.Fatalln("Error in GetAllUsers function")
// 	message := recover()
// 	log.Fatalln("ERROR: ", message)
// }

func GetAllUsers() ([]User, error) {
	rows, err := db.DB.Query("SELECT * FROM users WHERE name = some")
	if err != nil {
		return nil, err
	}
	// if rows != nil {
	// 	return nil, err
	// }
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err = rows.Scan(&u.ID, &u.Name, &u.Test)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
