package repository

import (
	defaultdto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"gorm.io/gorm"
)

type CookRepository interface {
	Create(entity *entity.Cook) error
	FindAll(pagination *defaultdto.PaginationRequest) ([]entity.Cook, int64, error)
	FindByID(id uint) (*entity.Cook, error)
	Update(entity *entity.Cook) error
	Delete(id uint) error
}

type cookRepository struct {
	db *gorm.DB
}

func NewCookRepository(db *gorm.DB) CookRepository {
	return &cookRepository{db: db}
}

func (r *cookRepository) Create(entity *entity.Cook) error {
	return r.db.Create(entity).Error
}

func (r *cookRepository) FindAll(pagination *defaultdto.PaginationRequest) ([]entity.Cook, int64, error) {
	var entities []entity.Cook
	var total int64

	query := r.db.Model(&entity.Cook{})

	
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

func (r *cookRepository) FindByID(id uint) (*entity.Cook, error) {
	var entity entity.Cook
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *cookRepository) Update(entity *entity.Cook) error {
	return r.db.Save(entity).Error
}

func (r *cookRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Cook{}, id).Error
}
