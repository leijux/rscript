package gui

import (
	"context"
	"log/slog"
	"net/netip"
	"slices"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/leijux/rscript/internal/pkg/engin"
	"github.com/leijux/rscript/internal/pkg/parser"
	"github.com/leijux/rscript/internal/pkg/version"
)

type App struct {
	path string

	ctx context.Context
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.path = "./default.yaml"
}

func (a *App) shutdown(ctx context.Context) {}

func (a *App) Test(engin.ProgressMsg, engin.ProgressResult) {}

func (a *App) GetConfigPath() string {
	return a.path
}

func (a *App) SetConfigPath() string {
	filepath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory: "./",
		DefaultFilename:  "default.yaml",
		Title:            "select a script",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "script",
				Pattern:     "*.yaml",
			},
		},
	})
	if err == nil && filepath != "" {
		a.path = filepath

		runtime.LogDebugf(a.ctx, "select script file %s", filepath)
		return filepath
	}
	return a.path
}

func (a *App) Start(clientIP string) bool {
	config := a.loadConfig()
	if config == nil {
		return false
	}
	if clientIP != "" {
		return a.retry(clientIP, config)
	}

	return a.run(config)
}

func (a *App) run(config *parser.RScript) bool {
	a.runScript(config)
	return true
}

func (a *App) retry(addrPortStr string, script *parser.RScript) bool {
	runtime.EventsEmit(a.ctx, "progress_info", engin.ProgressResult{
		ProgressMsg: engin.ProgressMsg{
			Name:    addrPortStr,
			Time:    time.Now().Format(time.DateTime),
			Msg:     "start retry",
			Percent: 0,
		},
	})

	addrPort, err := netip.ParseAddrPort(addrPortStr)
	if err != nil {
		runtime.EventsEmit(a.ctx, "err_msg", err.Error())
		return false
	}
	// find addr
	i := slices.IndexFunc(script.Remotes, func(r parser.Remote) bool {
		return r.AddrPort.Compare(addrPort) == 0
	})
	script.Remotes = []parser.Remote{script.Remotes[i]}

	runtime.LogInfof(a.ctx, "client ip %s retry", addrPortStr)
	a.runScript(script)

	return true
}

func (a *App) loadConfig() *parser.RScript {
	config, err := parser.ParseWithPath(a.path)
	if err != nil {
		runtime.EventsEmit(a.ctx, "err_msg", err.Error())
		return nil
	}
	return config
}

func (a *App) runScript(script *parser.RScript) {
	ctx, cancelFunc := context.WithCancel(a.ctx)
	e := engin.NewWithConfigAndContext(ctx, script, slog.Default())
	defer e.Close()

	cancelEventFunc := runtime.EventsOnce(a.ctx, "cancel_run", func(_ ...any) {
		cancelFunc()
	})
	defer cancelEventFunc()

	e.SetSendResult(func(result engin.ProgressResult) {
		runtime.EventsEmit(a.ctx, "progress_info", result)
	})

	e.Run()
}

func (a *App) GetVersion() string {
	return version.Version
}
