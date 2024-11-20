package tview

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type pager interface {
	Init(model) tea.Cmd
	Update(*model, tea.Msg) tea.Cmd
	View(model) string
}

type nextPageMsg struct{}

func nextPageAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return nextPageMsg{}
	})
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}
