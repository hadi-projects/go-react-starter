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
	if err := os.MkdirAll(cfg.Log.Dir, 0755); err != nil {
		panic(err)
	}

	SystemLogger = newLogger(*cfg, "system.log")
	AuthLogger = newLogger(*cfg, "auth.log")
	DBLogger = newLogger(*cfg, "db.log")
	RedisLogger = newLogger(*cfg, "redis.log")
	RateLimitLogger = newLogger(*cfg, "rate_limit.log")
}

func newLogger(cfg config.Config, fileName string) zerolog.Logger {
	// The provided diff changes the file logging mechanism from lumberjack to os.OpenFile
	// and the console output handling.
	// It also introduces a potential issue with `io.MultiWriter(writers...)` where `writers` is not defined.
	// Assuming the intent is to use `output` as the final writer for zerolog.New().

	// Original lumberjack setup:
	fileLogger := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Log.Dir, fileName), // Updated to cfg.Log.Dir
		MaxSize:    10,                                   // megabytes
		MaxBackups: 3,                                    // number of backups
		MaxAge:     28,                                   // days
		Compress:   true,
	}

	var writers []io.Writer
	writers = append(writers, fileLogger)

	if cfg.App.Env == "development" { // Updated to cfg.App.Env
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})
	}

	multi := io.MultiWriter(writers...)

	return zerolog.New(multi).With().Timestamp().Logger()
}
