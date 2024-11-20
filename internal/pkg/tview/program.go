package tview

import (
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Program struct {
	tp *tea.Program
}

func NewProgram(selectFileCh chan<- string) Program {
	return Program{
		tp: tea.NewProgram(NewModel(selectFileCh, selectFileCh != nil)),
	}
}

func (p Program) Run() {
	if _, err := p.tp.Run(); err != nil {
		slog.Error("program run err", "err", err)
		os.Exit(1)
	}
}

func (p Program) Send(msg tea.Msg) {
	p.tp.Send(msg)
}
