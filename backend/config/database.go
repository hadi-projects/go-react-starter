package config

type DatabaseConfig struct {
	Host         string
	Port         string
	UserName     string
	Password     string
	Name         string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  int
}
