package repository

import (
	"context"

	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role, permissionIDs []uint) error
	FindAll(ctx context.Context, pagination *dto.PaginationRequest) ([]entity.Role, int64, error)
	FindByID(ctx context.Context, id uint) (*entity.Role, error)
	FindByName(ctx context.Context, name string) (*entity.Role, error)
	Update(ctx context.Context, role *entity.Role, permissionIDs []uint) error
	Delete(ctx context.Context, id uint) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role, permissionIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err
		}

		if len(permissionIDs) > 0 {
			var permissions []entity.Permission
			if err := tx.Find(&permissions, permissionIDs).Error; err != nil {
				return err
			}
			if err := tx.Model(role).Association("Permissions").Replace(permissions); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *roleRepository) FindAll(ctx context.Context, pagination *dto.PaginationRequest) ([]entity.Role, int64, error) {
	var roles []entity.Role
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Role{})

	if pagination.Search != "" {
		searchTerm := "%" + pagination.Search + "%"
		query = query.Where("name LIKE ?", searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pagination.GetPage() - 1) * pagination.GetLimit()
	err := query.Order("id DESC").
		Preload("Permissions").
		Limit(pagination.GetLimit()).
		Offset(offset).
		Find(&roles).Error

	return roles, total, err
}

func (r *roleRepository) FindByID(ctx context.Context, id uint) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).Preload("Permissions").Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Update(ctx context.Context, role *entity.Role, permissionIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(role).Error; err != nil {
			return err
		}

		if len(permissionIDs) > 0 {
			var permissions []entity.Permission
			if err := tx.Find(&permissions, permissionIDs).Error; err != nil {
				return err
			}
			if err := tx.Model(role).Association("Permissions").Replace(permissions); err != nil {
				return err
			}
		} else {
			// If permissionIDs is empty, clear permissions? Or assume no change?
			// Usually in specific update, if it's passed as empty list, it might mean remove all.
			// Let's assume if it is explicitly passed (length 0 check might be checking nil slice vs empty slice, but here []uint is slice)
			// For safety, let's allow clearing if it's an empty slice, but the service logic should handle if it is nil.
			// However simple implementation: replace with provided list.
			if err := tx.Model(role).Association("Permissions").Clear(); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *roleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Clear associations first
		var role entity.Role
		if err := tx.First(&role, id).Error; err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Permissions").Clear(); err != nil {
			return err
		}
		// Delete role
		return tx.Delete(&role).Error
	})
}
