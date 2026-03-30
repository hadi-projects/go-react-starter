package repository

import (
	"context"

	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"gorm.io/gorm"
)

type ProdukRepository interface {
	Create(ctx context.Context, entity *entity.Produk) error
	FindAll(ctx context.Context, pagination *defaultDto.PaginationRequest) ([]entity.Produk, int64, error)
	FindByID(ctx context.Context, id uint) (*entity.Produk, error)
	Update(ctx context.Context, entity *entity.Produk) error
	Delete(ctx context.Context, id uint) error
}

type produkRepository struct {
	db *gorm.DB
}

func NewProdukRepository(db *gorm.DB) ProdukRepository {
	return &produkRepository{db: db}
}

func (r *produkRepository) Create(ctx context.Context, entity *entity.Produk) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *produkRepository) FindAll(ctx context.Context, pagination *defaultDto.PaginationRequest) ([]entity.Produk, int64, error) {
	var entities []entity.Produk
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Produk{})

	
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

func (r *produkRepository) FindByID(ctx context.Context, id uint) (*entity.Produk, error) {
	var entity entity.Produk
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *produkRepository) Update(ctx context.Context, entity *entity.Produk) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *produkRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Produk{}, id).Error
}
