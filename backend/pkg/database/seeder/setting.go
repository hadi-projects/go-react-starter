package seeder

import (
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"gorm.io/gorm"
)

func SeedSettings(db *gorm.DB) {
	settings := []entity.Setting{
		// Website Settings
		{Key: "app_name", Value: "Go-React Starter", Category: "website", FieldType: "text", Label: "Application Name", Description: "The name of your application shown in the browser title and sidebar."},
		{Key: "app_logo", Value: "", Category: "website", FieldType: "file", Label: "Logo", Description: "The logo image for your application."},
		{Key: "app_favicon", Value: "", Category: "website", FieldType: "file", Label: "Favicon", Description: "The icon shown in the browser tab."},

		// SMTP Settings
		{Key: "smtp_host", Value: "smtp.mailtrap.io", Category: "smtp", FieldType: "text", Label: "SMTP Host", Description: "The hostname of your SMTP server."},
		{Key: "smtp_port", Value: "2525", Category: "smtp", FieldType: "number", Label: "SMTP Port", Description: "The port of your SMTP server."},
		{Key: "smtp_user", Value: "", Category: "smtp", FieldType: "text", Label: "SMTP Username", Description: "The username for your SMTP server."},
		{Key: "smtp_pass", Value: "", Category: "smtp", FieldType: "password", Label: "SMTP Password", Description: "The password for your SMTP server."},
		{Key: "smtp_from_name", Value: "Admin Panel", Category: "smtp", FieldType: "text", Label: "From Name", Description: "The name that appears in the 'From' field of emails."},
		{Key: "smtp_from_email", Value: "no-reply@example.com", Category: "smtp", FieldType: "text", Label: "From Email", Description: "The email address used to send system notifications."},

		// Security Settings
		{Key: "jwt_secret", Value: "yoursecret", Category: "security", FieldType: "password", Label: "JWT Secret Key", Description: "The secret key used to sign and verify JSON Web Tokens."},
		{Key: "jwt_issuer", Value: "go-starter", Category: "security", FieldType: "text", Label: "JWT Issuer", Description: "The issuer name inside the JWT claims."},
		{Key: "jwt_access_expiration", Value: "15m", Category: "security", FieldType: "text", Label: "Access Token Expiration", Description: "How long before an access token expires (e.g., 15m, 1h)."},
		{Key: "jwt_refresh_expiration", Value: "7d", Category: "security", FieldType: "text", Label: "Refresh Token Expiration", Description: "How long before a refresh token expires (e.g., 24h, 7d)."},
		{Key: "cors_allowed_origins", Value: "http://localhost:5173", Category: "security", FieldType: "text", Label: "CORS Allowed Origins", Description: "Comma-separated list of allowed origins (e.g., http://localhost:5173, https://example.com)."},
		{Key: "rate_limit_rps", Value: "10", Category: "security", FieldType: "number", Label: "Rate Limit (RPS)", Description: "Number of requests allowed per second per client."},
		{Key: "rate_limit_burst", Value: "20", Category: "security", FieldType: "number", Label: "Rate Limit Burst", Description: "The maximum burst of requests allowed during a specific time interval."},

		// Storage Settings
		{Key: "storage_base_path", Value: "./storage/uploads", Category: "storage", FieldType: "text", Label: "Storage Base Path", Description: "The local directory where uploaded files are stored."},
		{Key: "storage_max_file_size_mb", Value: "50", Category: "storage", FieldType: "number", Label: "Max File Size (MB)", Description: "The maximum size for an individual uploaded file in Megabytes."},

		// Database Settings (Internal)
		{Key: "db_host", Value: "localhost", Category: "internal", FieldType: "text", Label: "Database Host", Description: "⚠️ Restart required. The hostname of your MySQL database."},
		{Key: "db_port", Value: "3306", Category: "internal", FieldType: "number", Label: "Database Port", Description: "⚠️ Restart required. The port of your MySQL database."},
		{Key: "db_username", Value: "root", Category: "internal", FieldType: "text", Label: "Database Username", Description: "⚠️ Restart required. The username for the MySQL database."},
		{Key: "db_password", Value: "", Category: "internal", FieldType: "password", Label: "Database Password", Description: "⚠️ Restart required. The password for the MySQL database."},
		{Key: "db_name", Value: "go_starter_db", Category: "internal", FieldType: "text", Label: "Database Name", Description: "⚠️ Restart required. The name of the MySQL database."},

		// Infrastructure Settings (Internal)
		{Key: "redis_host", Value: "localhost", Category: "internal", FieldType: "text", Label: "Redis Host", Description: "⚠️ Restart required. The hostname of your Redis server."},
		{Key: "redis_port", Value: "6379", Category: "internal", FieldType: "number", Label: "Redis Port", Description: "⚠️ Restart required. The port of your Redis server."},
		{Key: "redis_password", Value: "", Category: "internal", FieldType: "password", Label: "Redis Password", Description: "⚠️ Restart required. The password of your Redis server."},
		{Key: "kafka_brokers", Value: "localhost:9092", Category: "internal", FieldType: "text", Label: "Kafka Brokers", Description: "⚠️ Restart required. Comma-separated list of Kafka broker addresses."},
		{Key: "kafka_topic", Value: "password-reset", Category: "internal", FieldType: "text", Label: "Kafka Topic", Description: "⚠️ Restart required. The default topic for Kafka messages."},

		// Advance Settings
		{Key: "registration_open", Value: "true", Category: "advance", FieldType: "boolean", Label: "Registration Open", Description: "Whether new users can register themselves."},
		{Key: "maintenance_mode", Value: "false", Category: "advance", FieldType: "boolean", Label: "Maintenance Mode", Description: "Disable public access for maintenance."},
	}

	for _, s := range settings {
		var count int64
		db.Model(&entity.Setting{}).Where("`key` = ?", s.Key).Count(&count)
		if count == 0 {
			if err := db.Create(&s).Error; err != nil {
				logger.SystemLogger.Error().Err(err).Msgf("Failed to seed setting: %s", s.Key)
			}
		}
	}
	logger.SystemLogger.Info().Msg("Settings Seeding Completed!")
}
