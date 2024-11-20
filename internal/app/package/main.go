package main

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"net/netip"

	"github.com/alecthomas/kong"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"github.com/leijux/rscript/internal/pkg/engin"
	"github.com/leijux/rscript/internal/pkg/log"
	"github.com/leijux/rscript/internal/pkg/parser"
	"github.com/leijux/rscript/internal/pkg/tview"
)

//go:embed all:assets
var assets embed.FS

var CLI struct {
	Run struct {
		IP       []string `name:"ip"  short:"i" help:"remote ip" default:"192.168.0.1"`
		port     uint16   `name:"port" help:"remote port" default:"22"`
		Username string   `name:"username"  short:"u" help:"remote username" default:"root"`
		Password string   `name:"password"  short:"p" help:"remote password"`
	} `cmd:"" default:"withargs" help:"program run"`
}

// go build -ldflags "-s -w" -o upgrade.exe
func main() {
	ctx := kong.Parse(&CLI)
	if ctx.Error != nil {
		panic(ctx.Error)
	}
	// init log
	logger, lumberjackLogger := log.InitLog()
	defer lumberjackLogger.Close()

	//
	scriptFile, err := assets.ReadFile("assets/default.yaml")
	if err != nil {
		panic(err)
	}
	//
	script, err := parser.ParseWithBytes(scriptFile)
	if err != nil {
		panic(err)
	}

	//
	script.Remotes = script.Remotes[:0]
	for _, ip := range CLI.Run.IP {
		script.Remotes = append(script.Remotes, parser.Remote{
			AddrPort: netip.AddrPortFrom(netip.MustParseAddr(ip), CLI.Run.port),
			Username: CLI.Run.Username,
			Password: CLI.Run.Password,
		})
	}

	engin.SetInternalFunc("rscript.uploadFile", func(client *ssh.Client, args ...string) (string, error) {
		if len(args) != 2 {
			return "", errors.New("args error")
		}
		var (
			localFilePath  = args[0]
			remoteFilePath = args[1]
		)
		if localFilePath == "" {
			return "", errors.New("local file path is empty")
		}
		if remoteFilePath == "" {
			return "", errors.New("remote file path is empty")
		}

		sftpClient, err := sftp.NewClient(client)
		if err != nil {
			return "", fmt.Errorf("sftp.NewClient: %w", err)
		}
		defer sftpClient.Close()

		remoteFile, err := sftpClient.Create(remoteFilePath)
		if err != nil {
			return "", fmt.Errorf("sftpClient.Create: %w", err)
		}
		defer remoteFile.Close()

		packageFile, err := assets.Open(localFilePath)
		if err != nil {
			panic(err)
		}
		defer packageFile.Close()

		_, err = io.Copy(remoteFile, packageFile)
		if err != nil {
			return "", fmt.Errorf("io.Copy: %w", err)
		}
		return "", nil
	})

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
