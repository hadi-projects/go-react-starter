package repository

import (
	"github.com/hadi-projects/go-react-starter/internal/dto"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"gorm.io/gorm"
)

type AbcRepository interface {
	Create(entity *entity.Abc) error
	FindAll(pagination *dto.PaginationRequest) ([]entity.Abc, int64, error)
	FindByID(id uint) (*entity.Abc, error)
	Update(entity *entity.Abc) error
	Delete(id uint) error
}

type abcRepository struct {
	db *gorm.DB
}

func NewAbcRepository(db *gorm.DB) AbcRepository {
	return &abcRepository{db: db}
}

func (r *abcRepository) Create(entity *entity.Abc) error {
	return r.db.Create(entity).Error
}

func (r *abcRepository) FindAll(pagination *dto.PaginationRequest) ([]entity.Abc, int64, error) {
	var entities []entity.Abc
	var total int64

	query := r.db.Model(&entity.Abc{})

	if pagination.Search != "" {
		query = query.Where("name LIKE ?", "%"+pagination.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pagination.GetPage() - 1) * pagination.GetLimit()
	err := query.Limit(pagination.GetLimit()).
		Offset(offset).
		Find(&entities).Error

	return entities, total, err
}

func (r *abcRepository) FindByID(id uint) (*entity.Abc, error) {
	var entity entity.Abc
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *abcRepository) Update(entity *entity.Abc) error {
	return r.db.Save(entity).Error
}

func (r *abcRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Abc{}, id).Error
}
