package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/hadi-projects/go-react-starter/config"
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

func InitLogger(cfg *config.Config) {
	if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
		panic(err)
	}

	SystemLogger = newLogger(*cfg, "system.log")
	AuthLogger = newLogger(*cfg, "auth.log")
	DBLogger = newLogger(*cfg, "db.log")
	RedisLogger = newLogger(*cfg, "redis.log")
	RateLimitLogger = newLogger(*cfg, "rate_limit.log")
}

func newLogger(cfg config.Config, filename string) zerolog.Logger {
	fileLogger := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.LogDir, filename),
		MaxSize:    10, // megabytes
		MaxBackups: 3,  // number of backups
		MaxAge:     28, // days
		Compress:   true,
	}

	var writers []io.Writer
	writers = append(writers, fileLogger)

	if cfg.APPEnv == "development" {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})
	}

	multi := io.MultiWriter(writers...)

	return zerolog.New(multi).With().Timestamp().Logger()
}
