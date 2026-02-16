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
	AppPort string `mapstructure:"APP_PORT"`
	AppName string `mapstructure:"APP_NAME"`
	APPEnv  string `mapstructure:"APP_ENV"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUserName string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`

	DBMaxIdleConns int `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBMaxOpenConns int `mapstructure:"DB_MAX_OPEN_CONNS"`
	DBMaxLifetime  int `mapstructure:"DB_MAX_LIFETIME"`

	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	CORSAllowedOrigins   string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	CORSAllowedMethods   string `mapstructure:"CORS_ALLOWED_METHODS"`
	CORSAllowedHeaders   string `mapstructure:"CORS_ALLOWED_HEADERS"`
	CORSMaxAge           int    `mapstructure:"CORS_MAX_AGE"`
	CORSExposedHeaders   string `mapstructure:"CORS_EXPOSED_HEADERS"`
	CORSAllowCredentials bool   `mapstructure:"CORS_ALLOW_CREDENTIALS"`

	JwtSecret               string `mapstructure:"JWT_SECRET"`
	JwtIssuer               string `mapstructure:"JWT_ISSUER"`
	JwtAccessExpirationTime string `mapstructure:"JWT_ACCESS_EXPIRATION_TIME"`

	RateLimitRps   int `mapstructure:"RATE_LIMIT_RPS"`
	RateLimitBurst int `mapstructure:"RATE_LIMIT_BURST"`

	RequestTimeOut int `mapstructure:"REQUEST_TIMEOUT"`

	APIKey string `mapstructure:"API_KEY"`

	BCryptCost int `mapstructure:"BRCYPT_COST"`

	AdminEmail    string `mapstructure:"ADMIN_EMAIL"`
	AdminPassword string `mapstructure:"ADMIN_PASSWORD"`
	LogDir        string `mapstructure:"LOG_DIR"`
}

func LoadConfig() (config Config) {
	viper.SetDefault("LOG_DIR", "./storage/logs")
	viper.SetDefault("DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DB_MAX_OPEN_CONNS", 100)
	viper.SetDefault("DB_MAX_LIFETIME", 60) // minutes

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
		log.Printf("Warning: .env file not found, using system environtment variables")
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error mapping config: ", err)
	}
	return config
}
