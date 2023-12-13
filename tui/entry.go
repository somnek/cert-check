package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type entry struct {
	input    textinput.Model
	logs     string
	viewport viewport.Model
}

func (m entry) Init() tea.Cmd {
	return nil
}

func InitEntry() *entry {
	// text input
	t := textinput.New()
	t.Placeholder = "Type here..."
	t.Focus()
	t.CharLimit = 200
	t.Width = 200

	m := entry{
		input: t,
	}
	return &m
}

func (m entry) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.input, cmd = m.input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Tab):
			// notify Update on main model to display statusMsg
			return InitProject()
		}
	}
	return m, cmd
}

func (m entry) View() string {
	return m.input.View()
}