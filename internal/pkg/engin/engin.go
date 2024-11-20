package engin

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/leijux/rscript/internal/pkg/parser"
)

type Config struct {
	SSHConfigs []SSHConfig
	Commands   []string
}

func NewDefaultConfig(remotes []parser.Remote, commands []string) Config {
	var config Config

	for _, r := range remotes {
		config.SSHConfigs = append(config.SSHConfigs, SSHConfig{
			User:     r.Username,
			Password: r.Password,
			Host:     r.AddrPort.Addr().String(),
			Port:     r.AddrPort.Port(),

			addr: r.AddrPort,
		})
	}
	config.Commands = commands

	return config
}

type ProgressMsg struct {
	Name    string  `json:"name"`
	Msg     string  `json:"msg"`
	Percent float64 `json:"percent"`
	Time    string  `json:"time"`
}

type ProgressResult struct {
	ProgressMsg

	Err    string `json:"err"`
	Result string `json:"result"`
}

func (pr *ProgressResult) clearResult() {
	pr.Result = ""
	pr.Err = ""
}

type Engin struct {
	clients     []*ssh.Client
	sendResultF []func(msg ProgressResult)
	c           Config

	ctx context.Context
	l   *slog.Logger
}

func New(path string, l *slog.Logger) (*Engin, error) {
	config, err := parser.ParseWithPath(path)
	if err != nil {
		return nil, err
	}

	return NewWithConfig(config, l), nil
}

func NewWithConfig(script *parser.RScript, l *slog.Logger) *Engin {
	return NewWithConfigAndContext(context.Background(), script, l)
}

func NewWithConfigAndContext(ctx context.Context, config *parser.RScript, l *slog.Logger) *Engin {
	return &Engin{
		l:   l.With("server", "engin"),
		c:   NewDefaultConfig(config.Remotes, config.Commands),
		ctx: ctx,
	}
}

func (e *Engin) SetSendResult(sendResultF ...func(result ProgressResult)) {
	e.sendResultF = sendResultF
}

func (e *Engin) sendResult(pr ProgressResult) {
	pr.Time = time.Now().Format("2006-01-02 15:04:05")
	for _, f := range e.sendResultF {
		f(pr)
	}
}

func (e *Engin) connect() error {
	for _, sshConfig := range e.c.SSHConfigs {
		client, err := sshConfig.connect()
		if err != nil {
			e.l.Error("fail in ssh", "remote_ip", sshConfig.Host, "err", err)

			e.sendResult(ProgressResult{
				ProgressMsg: ProgressMsg{
					Name:    sshConfig.addr.String(),
					Percent: 0,
				},

				Err: err.Error(),
			})
			continue
		}
		e.clients = append(e.clients, client)
	}
	if len(e.clients) == 0 {
		return errors.New("fail in engin: no remote clients found")
	}
	return nil
}

func (e *Engin) Run() {
	select {
	case <-e.ctx.Done():
		return
	default:
		err := e.connect()
		if err != nil {
			e.l.Warn("fail in engin connect", "err", err)
			return
		}

		wg := &sync.WaitGroup{}
		for _, client := range e.clients {
			wg.Add(1)
			go func() {
				defer wg.Done()
				e.execCommands(client)
			}()
		}
		wg.Wait()
	}
}

func (e *Engin) Close() error {
	var errs error
	for _, client := range e.clients {
		errs = errors.Join(errs, client.Close())
	}
	return errs
}

func execInternalCommand(client *ssh.Client, command string) (string, error) {
	if client == nil {
		command = strings.Join([]string{"rscript.exec", command}, " ")
	}
	commands := strings.Split(strings.TrimSpace(command), " ")
	if iFunc, ok := internalFuncMap[commands[0]]; ok {
		output, err := iFunc(client, commands[1:]...)
		if err != nil {
			return "", err
		}
		return output, nil
	}
	return "", errors.New("fail in engin: unknown command")
}

// ExecuteCommand runs a command on the remote server
func execExternalCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	output, err := session.CombinedOutput(command)

	return string(output), err
}

func (e *Engin) execCommand(client *ssh.Client, command string) (string, error) {
	if strings.HasPrefix(command, prefix) {
		return execInternalCommand(client, command)
	}
	return execExternalCommand(client, command)
}

func (e *Engin) execCommands(client *ssh.Client) {
	var (
		remoteIP = client.RemoteAddr().String()
		cn       = len(e.c.Commands)
		l        = e.l.With("remote_ip", remoteIP)

		pr = ProgressResult{
			ProgressMsg: ProgressMsg{
				Name: remoteIP,
			},
		}
	)

	for i, command := range e.c.Commands {
		select {
		case <-e.ctx.Done():
			return
		default:
			pr.Msg = command
			pr.Percent = float64(i) / float64(cn)
			e.sendResult(pr)

			output, err := e.execCommand(client, command)
			if err != nil {
				l.Error("exec command err", "output", output, "command", command, "err", err)

				pr.Err = err.Error()
				e.sendResult(pr)
				return
			}

			if output == "" {
				output = "succeed"
			}

			pr.Result = output
			e.sendResult(pr)
			pr.clearResult()
		}
	}

	pr.Msg = "complete"
	pr.Result = pr.Msg
	pr.Percent = 1
	e.sendResult(pr)
}
