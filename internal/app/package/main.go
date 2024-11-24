package main

import (
	"embed"
	"net/netip"

	"github.com/alecthomas/kong"

	"github.com/leijux/rscript/internal/pkg/engin"
	"github.com/leijux/rscript/internal/pkg/log"
	"github.com/leijux/rscript/internal/pkg/parser"
	"github.com/leijux/rscript/internal/pkg/tview"
)

//go:embed all:assets
var assets embed.FS

var CLI struct {
	Run struct {
		IP       []string `name:"ip"        short:"i" help:"remote ip"       default:"192.168.0.1"`
		port     uint16   `name:"port"                help:"remote port"     default:"22"`
		Username string   `name:"username"  short:"u" help:"remote username" default:"root"`
		Password string   `name:"password"  short:"p" help:"remote password"`
	} `cmd:"" default:"withargs" help:"program run"`
}

// go build -ldflags "-s -w" -o example.exe
func main() {
	ctx := kong.Parse(&CLI)
	if ctx.Error != nil {
		panic(ctx.Error)
	}
	// init log
	logger, lumberjackLogger := log.InitLog()
	defer lumberjackLogger.Close()

	//从assets读取脚本文件
	scriptFile, err := assets.ReadFile("assets/default.yaml")
	if err != nil {
		panic(err)
	}

	//parse script file data
	script, err := parser.ParseWithBytes(scriptFile)
	if err != nil {
		panic(err)
	}

	script.Remotes = script.Remotes[:0]
	for _, ip := range CLI.Run.IP {
		script.Remotes = append(script.Remotes, parser.Remote{
			AddrPort: netip.AddrPortFrom(netip.MustParseAddr(ip), CLI.Run.port),
			Username: CLI.Run.Username,
			Password: CLI.Run.Password,
		})
	}

	// 重写rscript.uploadFile 从assets上传文件
	engin.SetInternalFunc("rscript.uploadFile", engin.UploadFileWithFs(assets))

	program := tview.NewProgram(nil)

	go func() {
		e := engin.NewWithConfig(script, logger)
		defer e.Close()

		e.SetSendResult(func(result engin.ProgressResult) {
			program.Send(result)
		})

		e.Run()
	}()

	program.Run()
}
