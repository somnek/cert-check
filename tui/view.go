package tui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	MARGIN           = 2
	PADDING          = 1
	COLOR_PINK       = lipgloss.Color("#E95678")
	COLOR_GRAY_1     = lipgloss.Color("#E3E4DB")
	COLOR_GRAY_2     = lipgloss.Color("#CDCDCD")
	COLOR_GRAY_3     = lipgloss.Color("#B8BABA")
	COLOR_GRAY_4     = lipgloss.Color("#626262")
	COLOR_DARK_GREEN = lipgloss.Color("#57886C")
)

var (
	styleSelected = lipgloss.NewStyle().Foreground(COLOR_PINK)
	styleNormal   = lipgloss.NewStyle().Foreground(COLOR_GRAY_1)
	styleTitle    = lipgloss.NewStyle().Foreground(COLOR_GRAY_2).Background(COLOR_DARK_GREEN).Bold(true)
	styleHelper1  = lipgloss.NewStyle().Foreground(COLOR_GRAY_4)
	styleHelper2  = lipgloss.NewStyle().Foreground(COLOR_GRAY_3)
)

func (m model) View() string {
	if !m.loaded {
		// to do, use spinner when loading
		return "🔭 dailing saved domains..."
	}
	return m.list.View()
}
