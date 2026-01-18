package session

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://192.168.1.2:5173", "http://binaryneedle.my.id:9111", "http://binaryneedle.my.id:9112", "http://localhost:8080", "http://192.168.1.2:8080", "https://subcommerce.mraihanz.my.id"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-CSRF-TOKEN"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
