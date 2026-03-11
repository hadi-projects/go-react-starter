package repository

import (
	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"gorm.io/gorm"
)

type ProdukRepository interface {
	Create(entity *entity.Produk) error
	FindAll(pagination *defaultDto.PaginationRequest) ([]entity.Produk, int64, error)
	FindByID(id uint) (*entity.Produk, error)
	Update(entity *entity.Produk) error
	Delete(id uint) error
}

type produkRepository struct {
	db *gorm.DB
}

func NewProdukRepository(db *gorm.DB) ProdukRepository {
	return &produkRepository{db: db}
}

func (r *produkRepository) Create(entity *entity.Produk) error {
	return r.db.Create(entity).Error
}

func (r *produkRepository) FindAll(pagination *defaultDto.PaginationRequest) ([]entity.Produk, int64, error) {
	var entities []entity.Produk
	var total int64

	query := r.db.Model(&entity.Produk{})

	
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

func (r *produkRepository) FindByID(id uint) (*entity.Produk, error) {
	var entity entity.Produk
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *produkRepository) Update(entity *entity.Produk) error {
	return r.db.Save(entity).Error
}

func (r *produkRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Produk{}, id).Error
}
