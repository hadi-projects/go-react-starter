package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/go-react-starter/config"
)

func CORSMiddleware(config *config.Config) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     strings.Split(config.CORS.AllowedOrigins, ","),
		AllowMethods:     strings.Split(config.CORS.AllowedMethods, ","),
		AllowHeaders:     strings.Split(config.CORS.AllowedHeaders, ","),
		ExposeHeaders:    strings.Split(config.CORS.ExposedHeaders, ","),
		AllowCredentials: config.CORS.AllowCredentials,
		MaxAge:           time.Duration(config.CORS.MaxAge) * time.Second,
	})
}
