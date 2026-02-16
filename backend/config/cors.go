package config

type CORSConfig struct {
	AllowedOrigins   string
	AllowedMethods   string
	AllowedHeaders   string
	MaxAge           int
	ExposedHeaders   string
	AllowCredentials bool
}
