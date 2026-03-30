package seeder

import (
	"context"
	"strconv"

	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	repository "github.com/hadi-projects/go-react-starter/internal/repository/default"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var users = []map[string]string{
	{"email": "admin@mail.com", "password": "Password@123", "roleId": "1"},
	{"email": "budi@mail.com", "password": "Password@123", "roleId": "2"},
	{"email": "annisa@mail.com", "password": "Password@123", "roleId": "2"},
}

func SeedUser(db *gorm.DB, bcryptCost int) {
	for _, userData := range users {
		email := userData["email"]
		password := userData["password"]
		roleIdStr := userData["roleId"]

		roleId, err := strconv.ParseUint(roleIdStr, 10, 32)
		if err != nil {
			logger.SystemLogger.Error().Err(err).Msgf("Failed to parse roleId for user %s", email)
			continue
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
		if err != nil {
			logger.SystemLogger.Error().Err(err).Msgf("Failed to hash password for user %s", email)
			continue
		}

		var user entity.User
		if err := db.Where("email = ?", email).First(&user).Error; err == nil {
			logger.SystemLogger.Info().Msgf("User %s already exists, skipping.", email)
			continue
		}

		user = entity.User{Email: email, Password: string(hashedPassword), RoleID: uint(roleId)}
		if err := repository.NewUserRepository(db).Create(context.Background(), &user); err != nil {
			logger.SystemLogger.Error().Err(err).Msgf("Failed to create user %s", email)
		} else {
			logger.SystemLogger.Info().Msgf("User %s created successfully.", email)
		}
	}

	logger.SystemLogger.Info().Int("bcrypt_cost", bcryptCost).Msg("User Seeding Completed!")
}
