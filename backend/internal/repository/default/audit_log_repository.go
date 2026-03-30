package repository

import (
	"context"

	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"gorm.io/gorm"
)

type AuditLogRepository interface {
	Create(ctx context.Context, log *logger.AuditLog) error
	FindAll(ctx context.Context, query *dto.AuditLogQuery) ([]entity.AuditLog, int64, error)
}

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, log *logger.AuditLog) error {
	entityLog := &entity.AuditLog{
		RequestID: log.RequestID,
		UserID:    log.UserID,
		UserEmail: log.UserEmail,
		Action:    log.Action,
		Module:    log.Module,
		TargetID:  log.TargetID,
		Metadata:  log.Metadata,
	}

	// Signal logger to skip this operation to avoid recursion if AuditLog uses GORM
	ctx = context.WithValue(ctx, logger.CtxKeySkipLogging, true)
	return r.db.WithContext(ctx).Create(entityLog).Error
}

func (r *auditLogRepository) FindAll(ctx context.Context, query *dto.AuditLogQuery) ([]entity.AuditLog, int64, error) {
	var logs []entity.AuditLog
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.AuditLog{})

	if query.Module != "" {
		db = db.Where("module = ?", query.Module)
	}
	if query.Action != "" {
		db = db.Where("action = ?", query.Action)
	}
	if query.UserEmail != "" {
		db = db.Where("user_email = ?", query.UserEmail)
	}
	if query.RequestID != "" {
		db = db.Where("request_id = ?", query.RequestID)
	}

	db.Count(&total)

	offset := (query.GetPage() - 1) * query.GetLimit()
	err := db.Order("id DESC").Offset(offset).Limit(query.GetLimit()).Find(&logs).Error

	return logs, total, err
}
