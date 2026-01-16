package db

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func errorInitDB() {
	log.Println("❌ Error in InitDB function")
	message := recover()
	log.Println("ERROR: ", message)
}

func InitDB(db_user string, db_pass string, db_host string, db_port string, db_name string) {
	dsn := "postgres://" + db_user + ":" + db_pass + "@" +
		db_host + ":" + db_port + "/" + db_name +
		"?sslmode=require&statement_cache_mode=describe"
	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		defer errorInitDB()
		errMessage := "❌ Failed to connect to database: " + err.Error()
		panic(errMessage)
	}

	// Connection pool settings (important for Supabase)
	DB.SetMaxOpenConns(5)
	DB.SetMaxIdleConns(2)
	DB.SetConnMaxLifetime(30 * time.Minute)

	if err = DB.Ping(); err != nil {
		log.Fatal("❌ Failed to ping DB:", err)
	}

	log.Println("✅ Database connected")
}
