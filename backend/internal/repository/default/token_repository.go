package repository

import (
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	"gorm.io/gorm"
)

type TokenRepository interface {
	Create(token *entity.PasswordResetToken) error
	FindByToken(token string) (*entity.PasswordResetToken, error)
	Delete(token *entity.PasswordResetToken) error
	DeleteByUserID(userID uint) error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Create(token *entity.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *tokenRepository) FindByToken(token string) (*entity.PasswordResetToken, error) {
	var t entity.PasswordResetToken
	err := r.db.Preload("User").Where("token = ?", token).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tokenRepository) Delete(token *entity.PasswordResetToken) error {
	return r.db.Delete(token).Error
}

func (r *tokenRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&entity.PasswordResetToken{}).Error
}
