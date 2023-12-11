package tui

import (
	"strings"

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
	var b strings.Builder
	title := "Cert Check"

	b.WriteString(styleTitle.Render(title))
	b.WriteString("\n\n")

	if m.page == 0 {
		// page 0: domain view

		// domain list
		for i, ssl := range m.ssls {
			var cursor string
			var style lipgloss.Style // style for each item

			if m.cursor == i {
				cursor = "│"
				style = styleSelected
			} else {
				cursor = " "
				style = styleNormal
			}

			writeDomainList(&b, ssl, style, cursor)
		}

		// print errors
		if m.err != nil {
			b.WriteString(m.err.Error() + "\n\n")
		}

		// logs
		b.WriteString(m.logs)
		b.WriteString("\n")

		writePageControls(&b, m.page)

	} else {
		// page 1: input view
		b.WriteString("Enter a domain name: \n\n")
		b.WriteString(m.input.View() + "\n\n")

		// print errors
		if m.err != nil {
			b.WriteString(m.err.Error() + "\n\n")
		}

		writePageControls(&b, m.page)
	}

	containerStyle := lipgloss.NewStyle().Margin(MARGIN).Padding(PADDING)
	return containerStyle.Render(b.String())
}
