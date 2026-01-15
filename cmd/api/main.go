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
	auth.POST("/register/seller", controller.CreateSellerHandler)
	auth.GET("/status", controller.CheckStatusHandler)

	// carts
	carts := v1.Group("/carts")
	carts.GET("/", controller.GetCartsHandler)
	carts.POST("/", controller.CreateCartHandler)
	carts.PATCH("/", controller.UpdateCartHandler)
	carts.DELETE("/", controller.DeleteCartsHandler)
	carts.DELETE("/product/:product_id/:product_variant_id", controller.DeleteCartHandler)

	// chats
	chats := v1.Group("/chats")
	// chats.GET("/", controller.GetAllConversationsHandler)
	chats.GET("/:id", controller.GetConversationHandler)
	// chats.GET("/messages/:id", controller.GetMessagesHandler)
	// chats.POST("/messages/:id", controller.CreateMessageHandler)

	// csrf
	csrfRoutes := v1.Group("/csrf")
	// generate csrf token
	csrfRoutes.GET("/", controller.CreateToken)
	// csrf token from session
	csrfRoutes.GET("/session", controller.GetToken)

	// order
	order := v1.Group("/orders")
	order.GET("/", controller.GetOrdersHandler)
	order.GET("/seller", controller.GetOrderSellerHandler)
	order.POST("/", controller.CreateOrderHandler)
	order.PATCH("/:order_id/:status_id", controller.UpdateStatusOrderHandler)
	order.GET("/checkouts", controller.GetCheckoutOrdersHandler)
	order.POST("/checkouts", controller.CreateCheckoutOrderHandler)
	order.GET("/payments", controller.GetOrderPaymentsHandler)
	// using product_id and product_variant_id in json
	order.PATCH("/", controller.UpdateOrderHandler)

	// products
	products := v1.Group("/products")
	products.GET("/", controller.GetProductsHandler)
	products.POST("/", controller.CreateProductHandler)
	products.POST("/:product_id/:product_variant_id", controller.CreateProductVariantsHandler)
	products.PATCH("/:product_id/:product_variant_id", controller.UpdateProductHandler)
	products.DELETE("/:product_id/:product_variant_id", controller.DeleteProductHandler)
	products.GET("/seller", controller.GetProductsSellerHandler)
	products.GET("/:id", controller.GetProductHandler)
	products.GET("/hot", controller.GetProductsHotHandler)
	products.GET("/discount", controller.GetProductsDiscountHandler)

	// ratings
	ratings := v1.Group("/ratings")
	ratings.GET("/:product_id", controller.GetRatingHandler)
	ratings.POST("/:product_id/:product_variant_id/:order_id", controller.CreateRatingHandler)
	ratings.GET("/comments/:product_id", controller.GetRatingCommentsHandler)

	// sellers
	sellers := v1.Group("/sellers")
	sellers.GET("/", controller.GetSeller)
	sellers.GET("/summarize/:product_id", controller.GetSellerSummarize)
	sellers.PATCH("/", controller.UpdateSellerHandler)
	sellers.DELETE("/", controller.DeleteSellerHandler)

	// users
	users := v1.Group("/users")
	users.GET("/", controller.GetUserHandler)
	users.PATCH("/", controller.UpdateUserHandler)
	users.DELETE("/", controller.DeleteUserHandler)

	// admins
	admins := v1.Group("/admins")
	// admin management by admin
	admins.GET("/", controller.GetAllAdminsHandler)
	admins.GET("/:name", controller.GetAdminHandler)
	admins.POST("/", controller.CreateAdminHandler)
	admins.PATCH("/:id", controller.UpdateAdminHandler)
	admins.DELETE("/:id", controller.DeleteAdminHandler)

	// user management by admin
	admins.GET("/users", controller.GetAllUsersHandler)
	admins.GET("/users/:name", controller.GetUserByAdminHandler)
	admins.POST("/users", controller.CreateUserByAdminHandler)
	admins.PATCH("/users/:id", controller.UpdateUserByAdminHandler)
	admins.DELETE("/users/:id", controller.DeleteUserByAdminHandler)

	// seller management by admin
	admins.GET("/sellers", controller.GetAllSellersHandler)
	admins.GET("/sellers/:name", controller.GetSellerByAdminHandler)
	admins.POST("/sellers/:uid", controller.CreateSellerByAdminHandler)
	admins.PATCH("/sellers/:id", controller.UpdateSellerByAdminHandler)
	admins.DELETE("/sellers/:id", controller.DeleteSellerByAdminHandler)

	// product management by admin
	admins.GET("/products", controller.GetAllProductsHandler)
	admins.GET("/products/:name", controller.GetProductsByAdminHandler)
	// admins.POST("/products", controller.CreateProductByAdminHandler)
	admins.PATCH("/products/:product_id/:active", controller.UpdateActiveProductHandler)
	admins.DELETE("/products/:product_id/:product_variant_id", controller.DeleteProductByAdminHandler)

	r.Run(os.Getenv("APP_PORT"))
}
