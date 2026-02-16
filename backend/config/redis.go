package config

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	TTL      int // Default cache TTL in seconds
}
