package repository

import (
	"context"

	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	"gorm.io/gorm"
)

type TokenRepository interface {
	Create(ctx context.Context, token *entity.PasswordResetToken) error
	FindByToken(ctx context.Context, token string) (*entity.PasswordResetToken, error)
	Delete(ctx context.Context, token *entity.PasswordResetToken) error
	DeleteByUserID(ctx context.Context, userID uint) error
	CreateRefreshToken(ctx context.Context, token *entity.RefreshToken) error
	FindByRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Create(ctx context.Context, token *entity.PasswordResetToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *tokenRepository) FindByToken(ctx context.Context, token string) (*entity.PasswordResetToken, error) {
	var t entity.PasswordResetToken
	err := r.db.WithContext(ctx).Preload("User").Where("token = ?", token).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tokenRepository) Delete(ctx context.Context, token *entity.PasswordResetToken) error {
	return r.db.WithContext(ctx).Delete(token).Error
}

func (r *tokenRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&entity.PasswordResetToken{}).Error
}

func (r *tokenRepository) CreateRefreshToken(ctx context.Context, token *entity.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *tokenRepository) FindByRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	var t entity.RefreshToken
	err := r.db.WithContext(ctx).Preload("User").Where("token = ?", token).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&entity.RefreshToken{}).Error
}
