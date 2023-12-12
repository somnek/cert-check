package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
		m.loaded = true
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
