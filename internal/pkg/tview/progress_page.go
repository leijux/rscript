package tview

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/leijux/rscript/internal/pkg/engin"
)

const (
	padding    = 2
	maxWidth   = 120
	defaultMsg = "......"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render

type progressPage struct {
	progress map[string]progressModelAndMsg
}

type progressModelAndMsg struct {
	model progress.Model
	msg   string
	fail  bool
}

func NewProgressPage() progressPage {
	return progressPage{
		make(map[string]progressModelAndMsg),
	}
}

var _ pager = (*progressPage)(nil)

func (p progressPage) Init(m model) tea.Cmd {
	return nil
}

func (p progressPage) Update(m *model, msg tea.Msg) tea.Cmd {
	defer func() {
		m.pages[m.pageIndex] = p
	}()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		for name, pg := range p.progress {
			pg.model.Width = msg.Width - padding*2 - 4
			if pg.model.Width > maxWidth {
				pg.model.Width = maxWidth
			}
			p.progress[name] = pg
		}
		return nil
	case engin.ProgressResult:
		var (
			cmd  tea.Cmd
			flag = true
		)
		pg, ok := p.progress[msg.Name]
		if !ok {
			pg = progressModelAndMsg{
				progress.New(progress.WithDefaultGradient()),
				defaultMsg,
				false,
			}
		}

		cmd = pg.model.SetPercent(msg.Percent)
		pg.msg = msg.Msg
		if msg.Err != "" {
			output := msg.Err + msg.Result
			if len(output) > 40 {
				output = output[:40]
			}
			pg.msg = output
			pg.fail = true
		}
		p.progress[msg.Name] = pg

		for _, pg := range p.progress {
			// 如果有没有完成就会卡住
			if 1.0-pg.model.Percent() > 0.001 {
				flag = false
			}
		}
		if flag && len(p.progress) != 0 {
			return tea.Batch(cmd, nextPageAfter(time.Second*5))
		}

		return cmd
	case progress.FrameMsg:
		var cmd tea.Cmd

		for name, pg := range p.progress {
			var progressModel tea.Model
			progressModel, cmd = pg.model.Update(msg)
			if cmd != nil {
				p.progress[name] = progressModelAndMsg{
					progressModel.(progress.Model),
					pg.msg,
					pg.fail,
				}
				break
			}
		}

		return cmd
	default:
		return nil
	}
}

func (p progressPage) View(m model) string {
	pad := strings.Repeat(" ", padding)

	var s strings.Builder

	s.WriteString("\n")

	keys := slices.Sorted(maps.Keys(p.progress))

	for _, name := range keys {
		pg := p.progress[name]

		s.WriteString(pad)
		s.WriteString(pg.model.View())

		s.WriteString("\n\n")

		s.WriteString(pad)

		var (
			styleFunc func(...string) string
			tips      string
		)

		if pg.fail {
			styleFunc = errorStyle
			tips = "执行异常"
		} else {
			styleFunc = helpStyle
			tips = "exec"
		}
		s.WriteString(styleFunc(fmt.Sprintf("[%s]%s: %s", name, tips, pg.msg)))

		s.WriteString("\n\n\n")
	}

	return s.String()
}
