package tview

import (
	"errors"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type scriptSelectPage struct {
	filepicker   filepicker.Model
	selectedFile string
	err          error
}

var _ pager = (*scriptSelectPage)(nil)

func NewScriptSelectPage(b bool) pager {
	if !b {
		return EmptyPage{}
	}
	fp := filepicker.New()
	fp.AllowedTypes = []string{".yaml"}
	fp.CurrentDirectory = "./"

	return scriptSelectPage{
		filepicker: fp,
	}
}

func (p scriptSelectPage) Init(m model) tea.Cmd {
	return p.filepicker.Init()
}

func (p scriptSelectPage) Update(m *model, msg tea.Msg) tea.Cmd {
	defer func() {
		m.pages[m.pageIndex] = p
	}()

	switch msg.(type) {
	case clearErrorMsg:
		p.err = nil
	}

	var (
		cmd      tea.Cmd
		isSelect bool
	)
	p.filepicker, cmd = p.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := p.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		p.selectedFile = path
		isSelect = true
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := p.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		p.err = errors.New(path + " Not a valid script")
		p.selectedFile = ""
		isSelect = false
		return tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	if isSelect && p.selectedFile != "" {
		m.selectFileCh <- p.selectedFile
		return nextPageAfter(time.Second)
	}

	return cmd
}

func (p scriptSelectPage) View(_ model) string {
	var s strings.Builder
	s.WriteString("\n  ")

	if p.err != nil {
		s.WriteString(p.filepicker.Styles.DisabledFile.Render(p.err.Error()))
	} else if p.selectedFile == "" {
		s.WriteString("please select a script:")
	} else {
		s.WriteString("script: " + p.filepicker.Styles.Selected.Render(p.selectedFile))
	}
	s.WriteString("\n\n" + p.filepicker.View() + "\n")

	return s.String()
}
