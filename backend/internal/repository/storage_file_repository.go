package repository

import (
	"context"

	defaultDto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"gorm.io/gorm"
)

type StorageFileRepository interface {
	Create(ctx context.Context, file *entity.StorageFile) error
	FindAll(ctx context.Context, userID uint, pagination *defaultDto.PaginationRequest) ([]entity.StorageFile, int64, error)
	FindByID(ctx context.Context, id uint) (*entity.StorageFile, error)
	FindByIDAndUserID(ctx context.Context, id, userID uint) (*entity.StorageFile, error)
	Update(ctx context.Context, file *entity.StorageFile) error
	Delete(ctx context.Context, id uint) error
	CountShareLinks(ctx context.Context, fileID uint) (int64, error)
	GetTotalSizeByUser(ctx context.Context, userID uint) (int64, error)
}

type storageFileRepository struct {
	db *gorm.DB
}

func NewStorageFileRepository(db *gorm.DB) StorageFileRepository {
	return &storageFileRepository{db: db}
}

func (r *storageFileRepository) Create(ctx context.Context, file *entity.StorageFile) error {
	return r.db.WithContext(ctx).Create(file).Error
}

func (r *storageFileRepository) FindAll(ctx context.Context, userID uint, pagination *defaultDto.PaginationRequest) ([]entity.StorageFile, int64, error) {
	var files []entity.StorageFile
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.StorageFile{}).Where("user_id = ?", userID)

	if pagination.Search != "" {
		query = query.Where("original_name LIKE ? OR description LIKE ?",
			"%"+pagination.Search+"%", "%"+pagination.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pagination.GetPage() - 1) * pagination.GetLimit()
	err := query.Order("id DESC").
		Limit(pagination.GetLimit()).
		Offset(offset).
		Find(&files).Error

	return files, total, err
}

func (r *storageFileRepository) FindByID(ctx context.Context, id uint) (*entity.StorageFile, error) {
	var file entity.StorageFile
	if err := r.db.WithContext(ctx).First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *storageFileRepository) FindByIDAndUserID(ctx context.Context, id, userID uint) (*entity.StorageFile, error) {
	var file entity.StorageFile
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *storageFileRepository) Update(ctx context.Context, file *entity.StorageFile) error {
	return r.db.WithContext(ctx).Save(file).Error
}

func (r *storageFileRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.StorageFile{}, id).Error
}

func (r *storageFileRepository) CountShareLinks(ctx context.Context, fileID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.ShareLink{}).
		Where("file_id = ? AND is_active = true", fileID).
		Count(&count).Error
	return count, err
}

func (r *storageFileRepository) GetTotalSizeByUser(ctx context.Context, userID uint) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&entity.StorageFile{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(size), 0)").
		Scan(&total).Error
	return total, err
}
