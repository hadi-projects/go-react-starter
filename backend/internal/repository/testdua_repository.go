package repository

import (
	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"gorm.io/gorm"
)

type TestduaRepository interface {
	Create(entity *entity.Testdua) error
	FindAll(pagination *defaultDto.PaginationRequest) ([]entity.Testdua, int64, error)
	FindByID(id uint) (*entity.Testdua, error)
	Update(entity *entity.Testdua) error
	Delete(id uint) error
}

type testduaRepository struct {
	db *gorm.DB
}

func NewTestduaRepository(db *gorm.DB) TestduaRepository {
	return &testduaRepository{db: db}
}

func (r *testduaRepository) Create(entity *entity.Testdua) error {
	return r.db.Create(entity).Error
}

func (r *testduaRepository) FindAll(pagination *defaultDto.PaginationRequest) ([]entity.Testdua, int64, error) {
	var entities []entity.Testdua
	var total int64

	query := r.db.Model(&entity.Testdua{})

	
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

func (r *testduaRepository) FindByID(id uint) (*entity.Testdua, error) {
	var entity entity.Testdua
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *testduaRepository) Update(entity *entity.Testdua) error {
	return r.db.Save(entity).Error
}

func (r *testduaRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Testdua{}, id).Error
}
