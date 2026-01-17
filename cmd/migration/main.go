package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func init() {
	var err error

	// Load .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") + "?sslmode=disable"
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		errMessage := "Failed to connect to database: " + err.Error()
		log.Fatalln(errMessage)
	}

	log.Println("Database connected")
}

func main() {
	var queryMigration []string

	users := `CREATE TABLE users 
	(id UUID PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	img TEXT,
	email TEXT NOT NULL,
	password VARCHAR(60) NOT NULL,
	created_at TIMESTAMP DEFAULT NOW() NOT NULL)`
	queryMigration = append(queryMigration, users)

	sellers := `CREATE TABLE sellers
	(id UUID PRIMARY KEY,
	user_id UUID CONSTRAINT fk_user_id
		REFERENCES users(id)
		ON DELETE CASCADE
		ON UPDATE CASCADE
		NOT NULL,
	name VARCHAR(50) NOT NULL,
	img TEXT,
	address TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT NOW() NOT NULL)`
	queryMigration = append(queryMigration, sellers)

	products := `CREATE TABLE products
	(id SERIAL PRIMARY KEY
	seller_id UUID CONSTRAINT fk_seller_id
		REFERENCES sellers(id)
		ON DELETE CASCADE
		ON UPDATE CASCADE
		NOT NULL,
	name VARCHAR(100) NOT NULL,
	img)`
	queryMigration = append(queryMigration, products)

	for _, v := range queryMigration {
		DB.Query(v)
	}
}
