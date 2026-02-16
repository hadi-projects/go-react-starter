package config

type JWTConfig struct {
	Secret               string
	Issuer               string
	AccessExpirationTime string
}
