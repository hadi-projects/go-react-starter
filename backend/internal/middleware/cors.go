package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/go-react-starter/config"
)

func CORS(config *config.Config) gin.HandlerFunc {
	corsRules := cors.DefaultConfig()

	corsRules.AllowOrigins = []string{config.CORSAllowedOrigins}
	corsRules.AllowMethods = []string{http.MethodPost}
	corsRules.AllowHeaders = []string{config.CORSAllowedHeaders}
	corsRules.ExposeHeaders = []string{config.CORSAllowedHeaders}
	corsRules.AllowCredentials = config.CORSAllowCredentials
	corsRules.MaxAge = time.Duration(config.CORSMaxAge * int(time.Hour))

	return cors.New(corsRules)
}
