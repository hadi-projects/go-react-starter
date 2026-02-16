package repository

import (
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(role *entity.Role, permissionIDs []uint) error
	FindAll() ([]entity.Role, error)
	FindByID(id uint) (*entity.Role, error)
	FindByName(name string) (*entity.Role, error)
	Update(role *entity.Role, permissionIDs []uint) error
	Delete(id uint) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(role *entity.Role, permissionIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
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

func (r *roleRepository) FindAll() ([]entity.Role, error) {
	var roles []entity.Role
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *roleRepository) FindByID(id uint) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindByName(name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Preload("Permissions").Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Update(role *entity.Role, permissionIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
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

func (r *roleRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
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
