package repository

import (
	"context"

	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	Create(ctx context.Context, permission *entity.Permission) error
	FindAll(ctx context.Context, pagination *dto.PaginationRequest) ([]entity.Permission, int64, error)
	FindByID(ctx context.Context, id uint) (*entity.Permission, error)
	FindByName(ctx context.Context, name string) (*entity.Permission, error)
	Update(ctx context.Context, permission *entity.Permission) error
	Delete(ctx context.Context, id uint) error
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(ctx context.Context, permission *entity.Permission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

func (r *permissionRepository) FindAll(ctx context.Context, pagination *dto.PaginationRequest) ([]entity.Permission, int64, error) {
	var permissions []entity.Permission
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Permission{})

	if pagination.Search != "" {
		searchTerm := "%" + pagination.Search + "%"
		query = query.Where("name LIKE ?", searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pagination.GetPage() - 1) * pagination.GetLimit()
	err := query.Order("id DESC").
		Limit(pagination.GetLimit()).
		Offset(offset).
		Find(&permissions).Error

	return permissions, total, err
}

func (r *permissionRepository) FindByID(ctx context.Context, id uint) (*entity.Permission, error) {
	var permission entity.Permission
	err := r.db.WithContext(ctx).First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) FindByName(ctx context.Context, name string) (*entity.Permission, error) {
	var permission entity.Permission
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) Update(ctx context.Context, permission *entity.Permission) error {
	return r.db.WithContext(ctx).Save(permission).Error
}

func (r *permissionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Permission{}, id).Error
}
