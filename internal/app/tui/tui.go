package tui

import (
	"log/slog"
	"os"

	"github.com/leijux/rscript/internal/pkg/engin"
	"github.com/leijux/rscript/internal/pkg/log"
	"github.com/leijux/rscript/internal/pkg/tview"
)

const defaultConfig = "./default.yaml"

func Main() {
	logger, lumberjackLogger := log.InitLog(slog.LevelError)
	defer lumberjackLogger.Close()

	var (
		program      tview.Program
		scriptPath   = ""
		selectFileCh chan string
	)

	_, err := os.Stat(defaultConfig)
	if os.IsNotExist(err) {
		selectFileCh = make(chan string)
		program = tview.NewProgram(selectFileCh)
	} else {
		scriptPath = defaultConfig
		program = tview.NewProgram(nil)
	}

	go func() {
		if scriptPath == "" {
			scriptPath = <-selectFileCh
		}

		e, err := engin.New(scriptPath, logger)
		if err != nil {
			panic(err)
		}
		defer e.Close()

		e.SetSendResult(func(result engin.ProgressResult) {
			program.Send(result)
		})

		e.Run()
	}()

	program.Run()
}
