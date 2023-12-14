package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := DocStyle.GetMargin()
		width := msg.Width - left - right
		height := msg.Height - top - bottom
		WindowSize.Height = height
		WindowSize.Width = width
		m.list.SetSize(width, height)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, keys.Tab):
			statusCmd := m.list.NewStatusMessage("You hit enter!")
			entry := InitEntry()
			model, cmd := entry.Update(WindowSize)
			return model, tea.Batch(cmd, statusCmd)
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
