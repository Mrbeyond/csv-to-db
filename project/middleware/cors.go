package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	AllowedOrigins := []string{
		"*",
	}
	return cors.New(
		cors.Config{
			AllowOrigins: AllowedOrigins,
			AllowMethods: []string{"PUT", "GET", "POST", "DELETE", "PATCH"},
			AllowHeaders: []string{"Origin", "Content-Length",
				"Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Cache-Control",
				"X-Requested-With",
			},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		},
	)
}
