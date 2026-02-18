package migrations

import (
	"log"

	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	"gorm.io/gorm"
)

func CreatePasswordResetTokensTable(db *gorm.DB) {
	err := db.AutoMigrate(&entity.PasswordResetToken{})
	if err != nil {
		log.Fatalf("Failed to migrate password_reset_tokens table: %v", err)
	}
}
