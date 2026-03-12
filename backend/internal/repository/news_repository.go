package repository

import (
	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"gorm.io/gorm"
)

type NewsRepository interface {
	Create(entity *entity.News) error
	FindAll(pagination *defaultDto.PaginationRequest) ([]entity.News, int64, error)
	FindByID(id uint) (*entity.News, error)
	Update(entity *entity.News) error
	Delete(id uint) error
}

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) NewsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) Create(entity *entity.News) error {
	return r.db.Create(entity).Error
}

func (r *newsRepository) FindAll(pagination *defaultDto.PaginationRequest) ([]entity.News, int64, error) {
	var entities []entity.News
	var total int64

	query := r.db.Model(&entity.News{})

	
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

func (r *newsRepository) FindByID(id uint) (*entity.News, error) {
	var entity entity.News
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *newsRepository) Update(entity *entity.News) error {
	return r.db.Save(entity).Error
}

func (r *newsRepository) Delete(id uint) error {
	return r.db.Delete(&entity.News{}, id).Error
}
