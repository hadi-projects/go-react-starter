package repository

import (
	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"gorm.io/gorm"
)

type BlogRepository interface {
	Create(entity *entity.Blog) error
	FindAll(pagination *defaultDto.PaginationRequest) ([]entity.Blog, int64, error)
	FindByID(id uint) (*entity.Blog, error)
	Update(entity *entity.Blog) error
	Delete(id uint) error
}

type blogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) BlogRepository {
	return &blogRepository{db: db}
}

func (r *blogRepository) Create(entity *entity.Blog) error {
	return r.db.Create(entity).Error
}

func (r *blogRepository) FindAll(pagination *defaultDto.PaginationRequest) ([]entity.Blog, int64, error) {
	var entities []entity.Blog
	var total int64

	query := r.db.Model(&entity.Blog{})

	
	if pagination.Search != "" {
		query = query.Where("name LIKE ?", "%"+pagination.Search+"%")
		query = query.Or("content LIKE ?", "%"+pagination.Search+"%")
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

func (r *blogRepository) FindByID(id uint) (*entity.Blog, error) {
	var entity entity.Blog
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *blogRepository) Update(entity *entity.Blog) error {
	return r.db.Save(entity).Error
}

func (r *blogRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Blog{}, id).Error
}
