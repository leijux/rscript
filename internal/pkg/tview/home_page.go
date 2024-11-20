package tview

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/leijux/rscript/internal/pkg/version"
)

type HomePage struct{}

var _ pager = (*HomePage)(nil)

func (h HomePage) Init(m model) tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("rscript-"+version.Version), nextPageAfter(0), tea.ClearScreen)
}

func (h HomePage) Update(m *model, msg tea.Msg) tea.Cmd {
	return nil
}

func (h HomePage) View(_ model) string {
	return ""
}
