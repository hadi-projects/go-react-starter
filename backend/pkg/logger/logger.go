package logger

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/hadi-projects/go-react-starter/config"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type contextKey string

const (
	CtxKeySkipLogging contextKey = "skip_logging"
	CtxKeyUserID      contextKey = "user_id"
	CtxKeyUserEmail   contextKey = "user_email"
	CtxKeyRequestID   contextKey = "request_id"
)

func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "... [truncated]"
}

type SystemLog struct {
	RequestID    string
	Method       string
	Path         string
	StatusCode   int
	Latency      int64
	RequestBody  string
	ResponseBody string
}

type SystemLogRepository interface {
	Create(ctx context.Context, log *SystemLog) error
}

type AuditLog struct {
	RequestID string
	UserID    uint
	UserEmail string
	Action    string
	Module    string
	TargetID  string
	Metadata  string
}

type AuditLogRepository interface {
	Create(ctx context.Context, log *AuditLog) error
}

func LogAudit(ctx context.Context, action, module, targetID, metadata string) {
	if AuditLogRepo == nil {
		return
	}

	userID := uint(0)
	userEmail := "system"
	requestID := ""

	// Extract from context (common keys used in Gin middleware)
	if val, ok := ctx.Value(CtxKeyUserID).(uint); ok {
		userID = val
	}
	if val, ok := ctx.Value(CtxKeyUserEmail).(string); ok {
		userEmail = val
	}
	if val, ok := ctx.Value(CtxKeyRequestID).(string); ok {
		requestID = val
	}

	// Truncate metadata to avoid oversized logs
	truncatedMetadata := Truncate(metadata, 65536)

	_ = AuditLogRepo.Create(ctx, &AuditLog{
		RequestID: requestID,
		UserID:    userID,
		UserEmail: userEmail,
		Action:    action,
		Module:    module,
		TargetID:  targetID,
		Metadata:  truncatedMetadata,
	})

	// Also log to file-based AuditLogger
	WithCtx(ctx, AuditLogger).Info().
		Str("action", action).
		Str("module", module).
		Str("target_id", targetID).
		Msg("audit operation")
}

func WithCtx(ctx context.Context, l zerolog.Logger) *zerolog.Logger {
	lc := l.With()

	if val, ok := ctx.Value(CtxKeyRequestID).(string); ok && val != "" {
		lc = lc.Str("request_id", val)
	}

	if val, ok := ctx.Value(CtxKeyUserID).(uint); ok && val != 0 {
		lc = lc.Uint("user_id", val)
	}

	if val, ok := ctx.Value(CtxKeyUserEmail).(string); ok && val != "" {
		lc = lc.Str("user_email", val)
	}

	logger := lc.Logger()
	return &logger
}

var (
	SystemLogRepo SystemLogRepository
	AuditLogRepo  AuditLogRepository
)

type Logger interface {
	Info() *zerolog.Event
	Error() *zerolog.Event
	Debug() *zerolog.Event
	Warn() *zerolog.Event
	Fatal() *zerolog.Event
}

type logger struct {
	zLog zerolog.Logger
}

func (l *logger) Info() *zerolog.Event {
	return l.zLog.Info()
}

func (l *logger) Error() *zerolog.Event {
	return l.zLog.Error()
}

func (l *logger) Debug() *zerolog.Event {
	return l.zLog.Debug()
}

func (l *logger) Warn() *zerolog.Event {
	return l.zLog.Warn()
}

func (l *logger) Fatal() *zerolog.Event {
	return l.zLog.Fatal()
}

// Global instances for backward compatibility if needed,
// strictly speaking we should move to DI for everything.
var (
	SystemLogger    zerolog.Logger
	DBLogger        zerolog.Logger
	RedisLogger     zerolog.Logger
	RateLimitLogger zerolog.Logger
	AuditLogger     zerolog.Logger
)

func InitLogger(cfg *config.Config) {
	if err := os.MkdirAll(cfg.Log.Dir, 0755); err != nil {
		panic(err)
	}

	SystemLogger = newZeroLogger(*cfg, "system.log")
	DBLogger = newZeroLogger(*cfg, "db.log")
	RedisLogger = newZeroLogger(*cfg, "redis.log")
	RateLimitLogger = newZeroLogger(*cfg, "rate_limit.log")
	AuditLogger = newZeroLogger(*cfg, "audit.log")
}

func NewLogger(cfg config.Config, fileName string) Logger {
	zLog := newZeroLogger(cfg, fileName)
	return &logger{zLog: zLog}
}

func newZeroLogger(cfg config.Config, fileName string) zerolog.Logger {
	fileLogger := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Log.Dir, fileName),
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	var writers []io.Writer
	writers = append(writers, fileLogger)

	if cfg.App.Env == "development" {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})
	}

	multi := io.MultiWriter(writers...)

	return zerolog.New(multi).With().Timestamp().Logger()
}
