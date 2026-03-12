package repository

import (
	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"gorm.io/gorm"
)

type MainnnRepository interface {
	Create(entity *entity.Mainnn) error
	FindAll(pagination *defaultDto.PaginationRequest) ([]entity.Mainnn, int64, error)
	FindByID(id uint) (*entity.Mainnn, error)
	Update(entity *entity.Mainnn) error
	Delete(id uint) error
}

type mainnnRepository struct {
	db *gorm.DB
}

func NewMainnnRepository(db *gorm.DB) MainnnRepository {
	return &mainnnRepository{db: db}
}

func (r *mainnnRepository) Create(entity *entity.Mainnn) error {
	return r.db.Create(entity).Error
}

func (r *mainnnRepository) FindAll(pagination *defaultDto.PaginationRequest) ([]entity.Mainnn, int64, error) {
	var entities []entity.Mainnn
	var total int64

	query := r.db.Model(&entity.Mainnn{})

	
	if pagination.Search != "" {
		query = query.Where("name LIKE ?", "%"+pagination.Search+"%")
		query = query.Or("makananan LIKE ?", "%"+pagination.Search+"%")
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

func (r *mainnnRepository) FindByID(id uint) (*entity.Mainnn, error) {
	var entity entity.Mainnn
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *mainnnRepository) Update(entity *entity.Mainnn) error {
	return r.db.Save(entity).Error
}

func (r *mainnnRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Mainnn{}, id).Error
}
