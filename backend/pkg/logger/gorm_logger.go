package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	Zerolog zerolog.Logger
}

func NewGormLogger(zLog zerolog.Logger) *GormLogger {
	return &GormLogger{Zerolog: zLog}
}

func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.Zerolog.Info().Msgf(msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Zerolog.Warn().Msgf(msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Zerolog.Error().Msgf(msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// Check if logging should be skipped for this context
	if skip, ok := ctx.Value(CtxKeySkipLogging).(bool); ok && skip {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// Truncate SQL for logging if it's too large (e.g., 64KB)
	truncatedSQL := Truncate(sql, 65536)

	status := 200
	if err != nil && err != gorm.ErrRecordNotFound {
		status = 500
	}

	event := l.Zerolog.Info()
	if err != nil && err != gorm.ErrRecordNotFound {
		event = l.Zerolog.Error().Err(err)
	}

	// Try to get request_id from context if available
	requestID := ""
	if rid, ok := ctx.Value(CtxKeyRequestID).(string); ok {
		requestID = rid
	}

	event.
		Str("request_id", requestID).
		Str("method", "DATABASE").
		Str("path", "mysql").
		Int("status_code", status).
		Int64("latency", elapsed.Milliseconds()).
		Str("request_body", truncatedSQL).
		Int64("rows_affected", rows).
		Msg("database operation")

	if SystemLogRepo != nil {
		_ = SystemLogRepo.Create(ctx, &SystemLog{
			RequestID:    requestID,
			Method:       "DATABASE",
			Path:         "mysql",
			StatusCode:   status,
			Latency:      elapsed.Milliseconds(),
			RequestBody:  truncatedSQL,
			ResponseBody: fmt.Sprintf("rows:%d", rows),
		})
	}
}
