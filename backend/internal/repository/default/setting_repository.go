package repository

import (
	"context"

	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	"gorm.io/gorm"
)

type SettingRepository interface {
	FindAll(ctx context.Context) ([]entity.Setting, error)
	FindByCategory(ctx context.Context, category string) ([]entity.Setting, error)
	FindByKey(ctx context.Context, key string) (*entity.Setting, error)
	Update(ctx context.Context, setting *entity.Setting) error
	BulkUpdate(ctx context.Context, settings []entity.Setting) error
}

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepository{db: db}
}

func (r *settingRepository) FindAll(ctx context.Context) ([]entity.Setting, error) {
	var settings []entity.Setting
	err := r.db.WithContext(ctx).Find(&settings).Error
	return settings, err
}

func (r *settingRepository) FindByCategory(ctx context.Context, category string) ([]entity.Setting, error) {
	var settings []entity.Setting
	err := r.db.WithContext(ctx).Where("category = ?", category).Find(&settings).Error
	return settings, err
}

func (r *settingRepository) FindByKey(ctx context.Context, key string) (*entity.Setting, error) {
	var setting entity.Setting
	err := r.db.WithContext(ctx).Where("`key` = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *settingRepository) Update(ctx context.Context, setting *entity.Setting) error {
	return r.db.WithContext(ctx).Save(setting).Error
}

func (r *settingRepository) BulkUpdate(ctx context.Context, settings []entity.Setting) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, s := range settings {
			// update only the value
			if err := tx.Model(&entity.Setting{}).Where("`key` = ?", s.Key).Update("value", s.Value).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
