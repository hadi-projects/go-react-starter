package repository

import (
	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"gorm.io/gorm"
)

type WisudaRepository interface {
	Create(entity *entity.Wisuda) error
	FindAll(pagination *defaultDto.PaginationRequest) ([]entity.Wisuda, int64, error)
	FindByID(id uint) (*entity.Wisuda, error)
	Update(entity *entity.Wisuda) error
	Delete(id uint) error
}

type wisudaRepository struct {
	db *gorm.DB
}

func NewWisudaRepository(db *gorm.DB) WisudaRepository {
	return &wisudaRepository{db: db}
}

func (r *wisudaRepository) Create(entity *entity.Wisuda) error {
	return r.db.Create(entity).Error
}

func (r *wisudaRepository) FindAll(pagination *defaultDto.PaginationRequest) ([]entity.Wisuda, int64, error) {
	var entities []entity.Wisuda
	var total int64

	query := r.db.Model(&entity.Wisuda{})

	
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

func (r *wisudaRepository) FindByID(id uint) (*entity.Wisuda, error) {
	var entity entity.Wisuda
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *wisudaRepository) Update(entity *entity.Wisuda) error {
	return r.db.Save(entity).Error
}

func (r *wisudaRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Wisuda{}, id).Error
}
