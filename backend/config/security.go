package config

type SecurityConfig struct {
	RequestTimeOut int
	APIKey         string
	BCryptCost     int
	AdminEmail     string
	AdminPassword  string
}
