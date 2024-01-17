package tui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	MARGIN           = 2
	PADDING          = 1
	COLOR_PINK       = lipgloss.Color("#E95678")
	COLOR_PINK_2     = lipgloss.Color("#bd5970")
	COLOR_GRAY_1     = lipgloss.Color("#E3E4DB")
	COLOR_GRAY_2     = lipgloss.Color("#CDCDCD")
	COLOR_GRAY_3     = lipgloss.Color("#B8BABA")
	COLOR_GRAY_4     = lipgloss.Color("#626262")
	COLOR_DARK_GREEN = lipgloss.Color("#57886C")
)

var (
	styleApp      = lipgloss.NewStyle().Padding(0, 2) // 1px top/bottom, 2px left/right
	DocStyle      = lipgloss.NewStyle().Margin(0, 2)
	stylePadEntry = lipgloss.NewStyle().Padding(0, 4)
	styleSelected = lipgloss.NewStyle().Foreground(COLOR_PINK)
	styleNormal   = lipgloss.NewStyle().Foreground(COLOR_GRAY_1)
	styleTitle    = lipgloss.NewStyle().
			Foreground(COLOR_GRAY_2).
			Bold(true)
	styleHelper1 = lipgloss.NewStyle().Foreground(COLOR_GRAY_4)
	styleHelper2 = lipgloss.NewStyle().Foreground(COLOR_GRAY_3)
)

func (m model) View() string {
	return styleApp.Render(m.list.View())
}
