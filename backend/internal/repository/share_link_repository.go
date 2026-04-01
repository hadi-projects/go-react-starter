package repository

import (
	"context"

	"github.com/hadi-projects/go-react-starter/internal/entity"
	"gorm.io/gorm"
)

type ShareLinkRepository interface {
	Create(ctx context.Context, link *entity.ShareLink) error
	FindByToken(ctx context.Context, token string) (*entity.ShareLink, error)
	FindByFileID(ctx context.Context, fileID uint) ([]entity.ShareLink, error)
	FindByID(ctx context.Context, id uint) (*entity.ShareLink, error)
	FindByIDAndUserID(ctx context.Context, id, userID uint) (*entity.ShareLink, error)
	Update(ctx context.Context, link *entity.ShareLink) error
	Delete(ctx context.Context, id uint) error
	RecordAccess(ctx context.Context, access *entity.ShareLinkAccess) error
	GetAccessLogs(ctx context.Context, shareLinkID uint) ([]entity.ShareLinkAccess, error)
}

type shareLinkRepository struct {
	db *gorm.DB
}

func NewShareLinkRepository(db *gorm.DB) ShareLinkRepository {
	return &shareLinkRepository{db: db}
}

func (r *shareLinkRepository) Create(ctx context.Context, link *entity.ShareLink) error {
	return r.db.WithContext(ctx).Create(link).Error
}

func (r *shareLinkRepository) FindByToken(ctx context.Context, token string) (*entity.ShareLink, error) {
	var link entity.ShareLink
	if err := r.db.WithContext(ctx).Preload("File").Where("token = ?", token).First(&link).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *shareLinkRepository) FindByFileID(ctx context.Context, fileID uint) ([]entity.ShareLink, error) {
	var links []entity.ShareLink
	err := r.db.WithContext(ctx).Where("file_id = ?", fileID).Order("id DESC").Find(&links).Error
	return links, err
}

func (r *shareLinkRepository) FindByID(ctx context.Context, id uint) (*entity.ShareLink, error) {
	var link entity.ShareLink
	if err := r.db.WithContext(ctx).First(&link, id).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *shareLinkRepository) FindByIDAndUserID(ctx context.Context, id, userID uint) (*entity.ShareLink, error) {
	var link entity.ShareLink
	err := r.db.WithContext(ctx).
		Joins("JOIN storage_files ON storage_files.id = share_links.file_id").
		Where("share_links.id = ? AND storage_files.user_id = ?", id, userID).
		First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *shareLinkRepository) Update(ctx context.Context, link *entity.ShareLink) error {
	return r.db.WithContext(ctx).Save(link).Error
}

func (r *shareLinkRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.ShareLink{}, id).Error
}

func (r *shareLinkRepository) RecordAccess(ctx context.Context, access *entity.ShareLinkAccess) error {
	return r.db.WithContext(ctx).Create(access).Error
}

func (r *shareLinkRepository) GetAccessLogs(ctx context.Context, shareLinkID uint) ([]entity.ShareLinkAccess, error) {
	var logs []entity.ShareLinkAccess
	err := r.db.WithContext(ctx).
		Where("share_link_id = ?", shareLinkID).
		Order("accessed_at DESC").
		Limit(100).
		Find(&logs).Error
	return logs, err
}
