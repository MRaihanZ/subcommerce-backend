package main

import (
	"log"
	"os"

	controller "github.com/MRaihanZ/subcommerce-backend/internal/controllers"
	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/session"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	// CORS (adjust origin as needed)
	r.Use(session.CorsMiddleware())

	// Session setup
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	db.InitDB(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	// Routes
	r.GET("/data", controller.GetUsers)
	r.Run(":8080")
	// r.POST("/data", controller.CreateUser)
	// r.POST("/login", controller.Login)
	// r.POST("/logout", controller.Logout)
}
