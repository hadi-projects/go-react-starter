package config

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Application struct {
	Config *Config
	Server *http.Server
	Router *gin.Engine
}

type Config struct {
	App       AppConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	CORS      CORSConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
	Security  SecurityConfig
	Log       LogConfig
}

func LoadConfig() (config Config) {
	viper.SetDefault("LOG_DIR", "./storage/logs")
	viper.SetDefault("DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DB_MAX_OPEN_CONNS", 100)
	viper.SetDefault("DB_MAX_LIFETIME", 60) // minutes
	viper.SetDefault("REDIS_TTL", 300)      // 5 minutes in seconds

	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	envVars := []string{
		"APP_PORT",
		"APP_NAME",
		"APP_ENV",
		"DB_HOST",
		"DB_PORT",
		"DB_USERNAME",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_MAX_IDLE_CONNS",
		"DB_MAX_OPEN_CONNS",
		"DB_MAX_LIFETIME",
		"REDIS_HOST",
		"REDIS_PORT",
		"REDIS_PASSWORD",
		"REDIS_DB",
		"REDIS_TTL",
		"CORS_ALLOWED_ORIGINS",
		"CORS_ALLOWED_METHODS",
		"CORS_ALLOWED_HEADERS",
		"CORS_MAX_AGE",
		"CORS_EXPOSED_HEADERS",
		"CORS_ALLOW_CREDENTIALS",
		"JWT_SECRET",
		"JWT_ISSUER",
		"JWT_ACCESS_EXPIRATION_TIME",
		"RATE_LIMIT_RPS",
		"RATE_LIMIT_BURST",
		"REQUEST_TIMEOUT",
		"API_KEY",
		"BCRYPT_COST",
		"ADMIN_EMAIL",
		"ADMIN_PASSWORD",
		"LOG_DIR",
	}

	for _, envVar := range envVars {
		viper.BindEnv(envVar)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	// Manually map to nested structs
	config.App = AppConfig{
		Port: viper.GetString("APP_PORT"),
		Name: viper.GetString("APP_NAME"),
		Env:  viper.GetString("APP_ENV"),
	}

	config.Database = DatabaseConfig{
		Host:         viper.GetString("DB_HOST"),
		Port:         viper.GetString("DB_PORT"),
		UserName:     viper.GetString("DB_USERNAME"),
		Password:     viper.GetString("DB_PASSWORD"),
		Name:         viper.GetString("DB_NAME"),
		MaxIdleConns: viper.GetInt("DB_MAX_IDLE_CONNS"),
		MaxOpenConns: viper.GetInt("DB_MAX_OPEN_CONNS"),
		MaxLifetime:  viper.GetInt("DB_MAX_LIFETIME"),
	}

	config.Redis = RedisConfig{
		Host:     viper.GetString("REDIS_HOST"),
		Port:     viper.GetString("REDIS_PORT"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       viper.GetInt("REDIS_DB"),
		TTL:      viper.GetInt("REDIS_TTL"),
	}

	config.CORS = CORSConfig{
		AllowedOrigins:   viper.GetString("CORS_ALLOWED_ORIGINS"),
		AllowedMethods:   viper.GetString("CORS_ALLOWED_METHODS"),
		AllowedHeaders:   viper.GetString("CORS_ALLOWED_HEADERS"),
		MaxAge:           viper.GetInt("CORS_MAX_AGE"),
		ExposedHeaders:   viper.GetString("CORS_EXPOSED_HEADERS"),
		AllowCredentials: viper.GetBool("CORS_ALLOW_CREDENTIALS"),
	}

	config.JWT = JWTConfig{
		Secret:               viper.GetString("JWT_SECRET"),
		Issuer:               viper.GetString("JWT_ISSUER"),
		AccessExpirationTime: viper.GetString("JWT_ACCESS_EXPIRATION_TIME"),
	}

	config.RateLimit = RateLimitConfig{
		Rps:   viper.GetInt("RATE_LIMIT_RPS"),
		Burst: viper.GetInt("RATE_LIMIT_BURST"),
	}

	config.Security = SecurityConfig{
		RequestTimeOut: viper.GetInt("REQUEST_TIMEOUT"),
		APIKey:         viper.GetString("API_KEY"),
		BCryptCost:     viper.GetInt("BCRYPT_COST"),
		AdminEmail:     viper.GetString("ADMIN_EMAIL"),
		AdminPassword:  viper.GetString("ADMIN_PASSWORD"),
	}

	config.Log = LogConfig{
		Dir: viper.GetString("LOG_DIR"),
	}

	return config
}
