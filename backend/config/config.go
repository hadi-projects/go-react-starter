package config

import (
	"fmt"
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
	Mail      MailConfig
	Kafka     KafkaConfig
	Frontend  FrontendConfig
	Storage   StorageConfig
}

type StorageConfig struct {
	BasePath      string
	MaxFileSizeMB int64
}

type MailConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	FromAddress string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type FrontendConfig struct {
	URL string
}

func LoadConfig() (config Config) {
	viper.SetDefault("LOG_DIR", "./storage/logs")
	viper.SetDefault("STORAGE_BASE_PATH", "./storage/uploads")
	viper.SetDefault("STORAGE_MAX_FILE_SIZE_MB", 50)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DB_MAX_OPEN_CONNS", 100)
	viper.SetDefault("DB_MAX_LIFETIME", 60) // minutes
	viper.SetDefault("REDIS_TTL", 300)      // 5 minutes in seconds
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_RETENTION_DAYS", 30)

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
		"JWT_REFRESH_EXPIRATION_TIME",
		"RATE_LIMIT_RPS",
		"RATE_LIMIT_BURST",
		"REQUEST_TIMEOUT",
		"API_KEY",
		"BCRYPT_COST",
		"ADMIN_EMAIL",
		"ADMIN_PASSWORD",
		"LOG_DIR",
		"LOG_LEVEL",
		"LOG_RETENTION_DAYS",
		"MAIL_HOST",
		"MAIL_PORT",
		"MAIL_USERNAME",
		"MAIL_PASSWORD",
		"MAIL_FROM_ADDRESS",
		"KAFKA_BROKERS",
		"KAFKA_TOPIC",
		"FRONTEND_URL",
		"STORAGE_BASE_PATH",
		"STORAGE_MAX_FILE_SIZE_MB",
	}

	for _, envVar := range envVars {
		viper.BindEnv(envVar)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
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
		Secret:                viper.GetString("JWT_SECRET"),
		Issuer:                viper.GetString("JWT_ISSUER"),
		AccessExpirationTime:  viper.GetString("JWT_ACCESS_EXPIRATION_TIME"),
		RefreshExpirationTime: viper.GetString("JWT_REFRESH_EXPIRATION_TIME"),
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
		Dir:           viper.GetString("LOG_DIR"),
		Level:         viper.GetString("LOG_LEVEL"),
		RetentionDays: viper.GetInt("LOG_RETENTION_DAYS"),
	}

	config.Mail = MailConfig{
		Host:        viper.GetString("MAIL_HOST"),
		Port:        viper.GetInt("MAIL_PORT"),
		User:        viper.GetString("MAIL_USERNAME"),
		Password:    viper.GetString("MAIL_PASSWORD"),
		FromAddress: viper.GetString("MAIL_FROM_ADDRESS"),
	}

	config.Kafka = KafkaConfig{
		Brokers: viper.GetStringSlice("KAFKA_BROKERS"),
		Topic:   viper.GetString("KAFKA_TOPIC"),
	}

	config.Frontend = FrontendConfig{
		URL: viper.GetString("FRONTEND_URL"),
	}

	config.Storage = StorageConfig{
		BasePath:      viper.GetString("STORAGE_BASE_PATH"),
		MaxFileSizeMB: viper.GetInt64("STORAGE_MAX_FILE_SIZE_MB"),
	}

	return config
}
