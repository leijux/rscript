package log

import (
	"log/slog"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/leijux/rscript/internal/pkg/version"
)

func InitLog(level slog.Level) (*slog.Logger, *lumberjack.Logger) {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./logs/rscript.log",
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     30,   //days
		Compress:   true, // disabled by default
	}

	zerologLogger := zerolog.New(lumberjackLogger)
	logger := slog.New(slogzerolog.Option{Level: level, Logger: &zerologLogger}.NewZerologHandler()).
		With("version", version.Version)

	slog.SetDefault(logger)

	return logger, lumberjackLogger
}
