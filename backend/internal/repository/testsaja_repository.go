package repository

import (
	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"gorm.io/gorm"
)

type TestsajaRepository interface {
	Create(entity *entity.Testsaja) error
	FindAll(pagination *defaultDto.PaginationRequest) ([]entity.Testsaja, int64, error)
	FindByID(id uint) (*entity.Testsaja, error)
	Update(entity *entity.Testsaja) error
	Delete(id uint) error
}

type testsajaRepository struct {
	db *gorm.DB
}

func NewTestsajaRepository(db *gorm.DB) TestsajaRepository {
	return &testsajaRepository{db: db}
}

func (r *testsajaRepository) Create(entity *entity.Testsaja) error {
	return r.db.Create(entity).Error
}

func (r *testsajaRepository) FindAll(pagination *defaultDto.PaginationRequest) ([]entity.Testsaja, int64, error) {
	var entities []entity.Testsaja
	var total int64

	query := r.db.Model(&entity.Testsaja{})

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

func (r *testsajaRepository) FindByID(id uint) (*entity.Testsaja, error) {
	var entity entity.Testsaja
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *testsajaRepository) Update(entity *entity.Testsaja) error {
	return r.db.Save(entity).Error
}

func (r *testsajaRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Testsaja{}, id).Error
}
