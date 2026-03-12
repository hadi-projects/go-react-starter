package repository

import (
	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"gorm.io/gorm"
)

type ArsipRepository interface {
	Create(entity *entity.Arsip) error
	FindAll(pagination *defaultDto.PaginationRequest) ([]entity.Arsip, int64, error)
	FindByID(id uint) (*entity.Arsip, error)
	Update(entity *entity.Arsip) error
	Delete(id uint) error
}

type arsipRepository struct {
	db *gorm.DB
}

func NewArsipRepository(db *gorm.DB) ArsipRepository {
	return &arsipRepository{db: db}
}

func (r *arsipRepository) Create(entity *entity.Arsip) error {
	return r.db.Create(entity).Error
}

func (r *arsipRepository) FindAll(pagination *defaultDto.PaginationRequest) ([]entity.Arsip, int64, error) {
	var entities []entity.Arsip
	var total int64

	query := r.db.Model(&entity.Arsip{})

	
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

func (r *arsipRepository) FindByID(id uint) (*entity.Arsip, error) {
	var entity entity.Arsip
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *arsipRepository) Update(entity *entity.Arsip) error {
	return r.db.Save(entity).Error
}

func (r *arsipRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Arsip{}, id).Error
}
