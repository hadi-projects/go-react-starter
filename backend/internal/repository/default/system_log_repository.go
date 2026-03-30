package repository

import (
	"context"

	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"gorm.io/gorm"
)

type SystemLogRepository interface {
	Create(ctx context.Context, log *logger.SystemLog) error
	FindAll(ctx context.Context, query *dto.SystemLogQuery) ([]entity.SystemLog, int64, error)
}

type systemLogRepository struct {
	db *gorm.DB
}

func NewSystemLogRepository(db *gorm.DB) SystemLogRepository {
	return &systemLogRepository{db: db}
}

func (r *systemLogRepository) Create(ctx context.Context, log *logger.SystemLog) error {
	entityLog := &entity.SystemLog{
		RequestID:    log.RequestID,
		Method:       log.Method,
		Path:         log.Path,
		StatusCode:   log.StatusCode,
		Latency:      log.Latency,
		RequestBody:  log.RequestBody,
		ResponseBody: log.ResponseBody,
	}
	
	// Use a context that signals the logger to skip this operation
	ctx = context.WithValue(ctx, logger.CtxKeySkipLogging, true)
	return r.db.WithContext(ctx).Create(entityLog).Error
}

func (r *systemLogRepository) FindAll(ctx context.Context, query *dto.SystemLogQuery) ([]entity.SystemLog, int64, error) {
	var logs []entity.SystemLog
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.SystemLog{})

	if query.Method != "" {
		db = db.Where("method = ?", query.Method)
	}
	if query.StatusCode != 0 {
		db = db.Where("status_code = ?", query.StatusCode)
	}
	if query.Path != "" {
		db = db.Where("path LIKE ?", "%"+query.Path+"%")
	}
	if query.RequestID != "" {
		db = db.Where("request_id = ?", query.RequestID)
	}

	db.Count(&total)

	offset := (query.GetPage() - 1) * query.GetLimit()
	err := db.Order("id DESC").Offset(offset).Limit(query.GetLimit()).Find(&logs).Error

	return logs, total, err
}
