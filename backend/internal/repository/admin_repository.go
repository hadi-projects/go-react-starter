package repository

import (
	defaultdto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"gorm.io/gorm"
)

type AdminRepository interface {
	Create(entity *entity.Admin) error
	FindAll(pagination *defaultdto.PaginationRequest) ([]entity.Admin, int64, error)
	FindByID(id uint) (*entity.Admin, error)
	Update(entity *entity.Admin) error
	Delete(id uint) error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) Create(entity *entity.Admin) error {
	return r.db.Create(entity).Error
}

func (r *adminRepository) FindAll(pagination *defaultdto.PaginationRequest) ([]entity.Admin, int64, error) {
	var entities []entity.Admin
	var total int64

	query := r.db.Model(&entity.Admin{})

	
	if pagination.Search != "" {
		query = query.Where("name LIKE ?", "%"+pagination.Search+"%")
	}
	

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pagination.GetPage() - 1) * pagination.GetLimit()
	err := query.Order("id DESC").
		Limit(pagination.GetLimit()).
		Offset(offset).
		Find(&entities).Error

	return entities, total, err
}

func (r *adminRepository) FindByID(id uint) (*entity.Admin, error) {
	var entity entity.Admin
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *adminRepository) Update(entity *entity.Admin) error {
	return r.db.Save(entity).Error
}

func (r *adminRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Admin{}, id).Error
}
