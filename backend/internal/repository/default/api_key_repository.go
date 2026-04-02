package repository

import (
	"context"
	"time"

	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	"gorm.io/gorm"
)

type ApiKeyRepository interface {
	Create(ctx context.Context, apiKey *entity.ApiKey) error
	FindAll(ctx context.Context, pagination *dto.PaginationRequest) ([]entity.ApiKey, int64, error)
	FindByID(ctx context.Context, id uint) (*entity.ApiKey, error)
	FindByHash(ctx context.Context, hash string) (*entity.ApiKey, error)
	Delete(ctx context.Context, id uint) error
	UpdateLastUsed(ctx context.Context, id uint) error
}

type apiKeyRepository struct {
	db *gorm.DB
}

func NewApiKeyRepository(db *gorm.DB) ApiKeyRepository {
	return &apiKeyRepository{db: db}
}

func (r *apiKeyRepository) Create(ctx context.Context, apiKey *entity.ApiKey) error {
	return r.db.WithContext(ctx).Create(apiKey).Error
}

func (r *apiKeyRepository) FindAll(ctx context.Context, pagination *dto.PaginationRequest) ([]entity.ApiKey, int64, error) {
	var keys []entity.ApiKey
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.ApiKey{})

	if pagination.Search != "" {
		searchTerm := "%" + pagination.Search + "%"
		query = query.Where("name LIKE ? OR prefix LIKE ?", searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pagination.GetPage() - 1) * pagination.GetLimit()
	err := query.Order("id DESC").
		Preload("Role").
		Limit(pagination.GetLimit()).
		Offset(offset).
		Find(&keys).Error

	return keys, total, err
}

func (r *apiKeyRepository) FindByID(ctx context.Context, id uint) (*entity.ApiKey, error) {
	var key entity.ApiKey
	err := r.db.WithContext(ctx).Preload("Role.Permissions").First(&key, id).Error
	return &key, err
}

func (r *apiKeyRepository) FindByHash(ctx context.Context, hash string) (*entity.ApiKey, error) {
	var key entity.ApiKey
	err := r.db.WithContext(ctx).
		Preload("Role.Permissions").
		Where("key_hash = ?", hash).
		First(&key).Error
	return &key, err
}

func (r *apiKeyRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.ApiKey{}, id).Error
}

func (r *apiKeyRepository) UpdateLastUsed(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.ApiKey{}).Where("id = ?", id).Update("last_used_at", &now).Error
}
