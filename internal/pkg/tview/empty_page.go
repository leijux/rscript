package tview

import tea "github.com/charmbracelet/bubbletea"

type EmptyPage struct{}

var _ pager = (*EmptyPage)(nil)

func (e EmptyPage) Init(m model) tea.Cmd {
	return nextPageAfter(0)
}

func (e EmptyPage) Update(_ *model, _ tea.Msg) tea.Cmd {
	return nil
}

func (e EmptyPage) View(_ model) string {
	return ""
}
