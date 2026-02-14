package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	SystemLogger    zerolog.Logger
	AuthLogger      zerolog.Logger
	DBLogger        zerolog.Logger
	RedisLogger     zerolog.Logger
	RateLimitLogger zerolog.Logger
)

type Config struct {
	LogDir       string
	Environtment string
}

func InitLogger(cfg Config) {
	if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
		panic(err)
	}

	SystemLogger = newLogger(cfg, "system.log")
	AuthLogger = newLogger(cfg, "auth.log")
	DBLogger = newLogger(cfg, "db.log")
	RedisLogger = newLogger(cfg, "redis.log")
	RateLimitLogger = newLogger(cfg, "rate_limit.log")
}

func newLogger(cfg Config, filename string) zerolog.Logger {
	fileLogger := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.LogDir, filename),
		MaxSize:    10, // megabytes
		MaxBackups: 3,  // number of backups
		MaxAge:     28, // days
		Compress:   true,
	}

	var writers []io.Writer
	writers = append(writers, fileLogger)

	if cfg.Environtment == "development" {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})
	}

	multi := io.MultiWriter(writers...)

	return zerolog.New(multi).With().Timestamp().Logger()
}
