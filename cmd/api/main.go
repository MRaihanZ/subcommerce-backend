package main

import (
	"log"
	"os"

	"github.com/MRaihanZ/subcommerce-backend/internal/controller"
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

	// Session setup
	store := cookie.NewStore([]byte("secret"))

	// (CORS, Session)
	r.Use(session.CorsMiddleware(), sessions.Sessions("mysession", store))

	db.InitDB(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	// Routes
	// Auth
	r.POST("/api/v1/auth/register", controller.CreateUserHandler)
	r.POST("/api/v1/auth/login", controller.VerifyUserHandler)

	// users
	r.GET("/api/v1/users", controller.GetUsersHandler)
	r.GET("/api/v1/users/:id", controller.GetUserHandler)

	// products
	// r.GET("/api/v1/products", controller.GetUser)
	r.Run(":8080")
}
