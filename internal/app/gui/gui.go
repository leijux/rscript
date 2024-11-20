package gui

import (
	"embed"
	"log/slog"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/leijux/rscript/internal/pkg/log"
	"github.com/leijux/rscript/internal/pkg/version"
)

func Main(assets embed.FS) error {
	logger, lumberjackLogger := log.InitLog()
	defer lumberjackLogger.Close()

	app := &App{}

	return wails.Run(&options.App{
		Title:            "rscript_" + version.Version,
		Width:            1024,
		Height:           768,
		MinWidth:         1024,
		MinHeight:        500,
		DisableResize:    false,
		AssetServer:      &assetserver.Options{Assets: assets},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind:             []any{app},
		Logger:           WailsLogger{Logger: logger.With("server", "wails")},
	})
}

type WailsLogger struct {
	*slog.Logger
}

func (l WailsLogger) Print(message string) {
	l.Logger.Debug(message)
}

func (l WailsLogger) Trace(message string) {
	l.Logger.Debug(message)
}

func (l WailsLogger) Debug(message string) {
	l.Logger.Debug(message)
}

func (l WailsLogger) Info(message string) {
	l.Logger.Info(message)
}

func (l WailsLogger) Warning(message string) {
	l.Logger.Warn(message)
}

func (l WailsLogger) Error(message string) {
	l.Logger.Error(message)
}

func (l WailsLogger) Fatal(message string) {
	l.Logger.Error(message)
	os.Exit(1)
}
