package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/MRaihanZ/subcommerce-backend/internal/controller"
	"github.com/MRaihanZ/subcommerce-backend/internal/db"
	"github.com/MRaihanZ/subcommerce-backend/internal/session"
	csrf "github.com/utrack/gin-csrf"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func decodeBase64(str string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Fatalf("Failed to decode base64 secret: %v", err)
	}
	return decoded
}

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	// session setup
	store := memstore.NewStore(decodeBase64(os.Getenv("SESSION_SECRET_CURRENT")), decodeBase64(os.Getenv("SESSION_SECRET_OLD")))
	// set cookie options
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	// (CORS, session)
	r.Use(session.CorsMiddleware(), sessions.Sessions("session_id", store))

	// CSRF middleware
	r.Use(csrf.Middleware(csrf.Options{
		Secret: string(decodeBase64(os.Getenv("CSRF_SECRET"))),
		ErrorFunc: func(c *gin.Context) {
			c.JSON(400, gin.H{"error": "CSRF token mismatch"})
			c.Abort()
		},
	}))

	db.InitDB(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	// routes
	v1 := r.Group("/api/v1")

	// auth
	auth := v1.Group("/auth")
	auth.POST("/login", controller.VerifyUserHandler)
	auth.POST("/logout", controller.LogoutHandler)
	auth.POST("/register", controller.CreateUserHandler)
	// auth.POST("/verification/seller", controller.VerificationSellerHandler)
	// auth.POST("/status/seller", controller.CheckStatusSellerHandler)
	// auth.POST("/register/seller", controller.CreateSellerHandler)
	auth.GET("/status", controller.CheckStatusHandler)

	// carts
	carts := v1.Group("/carts")
	carts.GET("/", controller.GetCartsHandler)
	carts.POST("/", controller.CreateCartHandler)
	carts.PATCH("/", controller.UpdateCartHandler)
	carts.DELETE("/", controller.DeleteCartsHandler)
	carts.DELETE("/product", controller.DeleteCartHandler)

	// csrf
	csrfRoutes := v1.Group("/csrf")
	// generate csrf token
	csrfRoutes.GET("/", controller.CreateToken)
	// csrf token from session
	csrfRoutes.GET("/session", controller.GetToken)

	// products
	products := v1.Group("/products")
	products.GET("/", controller.GetProductsHandler)
	products.GET("/:id", controller.GetProductHandler)
	products.GET("/hot", controller.GetProductsHotHandler)
	products.GET("/discount", controller.GetProductsDiscountHandler)

	// ratings
	ratings := v1.Group("/ratings")
	ratings.GET("/:product_id", controller.GetRatingHandler)
	ratings.GET("/comments/:product_id", controller.GetRatingCommentsHandler)

	// sellers
	sellers := v1.Group("/sellers")
	sellers.GET("/", controller.GetSeller)
	sellers.GET("/:product_id/summarize", controller.GetSellerSummarize)

	// users
	users := v1.Group("/users")
	users.GET("/", controller.GetUserHandler)
	users.PATCH("/", controller.UpdateUserHandler)
	// users.DELETE("/", controller.DeleteUserHandler)

	r.Run(os.Getenv("APP_PORT"))
}
