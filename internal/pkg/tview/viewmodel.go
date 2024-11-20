package tview

import (
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	pages     []pager
	pageIndex uint

	selectFileCh chan<- string
}

func NewModel(selectFileCh chan<- string, selectFile bool) model {
	return model{
		selectFileCh: selectFileCh,
		pages: []pager{
			HomePage{},
			NewScriptSelectPage(selectFile),
			NewProgressPage(),
			EmptyPage{},
		},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.pages[0].Init(m))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case nextPageMsg:
		if int(m.pageIndex) < len(m.pages)-1 {
			m.pageIndex++
			return m, m.pages[m.pageIndex].Init(m)
		}
		return m, tea.Quit
	}

	return m, m.pages[m.pageIndex].Update(&m, msg)
}

func (m model) View() string {
	return m.pages[m.pageIndex].View(m)
}
