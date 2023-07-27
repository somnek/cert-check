package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	inputs   []textinput.Model
}

func initialMode() model {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	return model{
		choices:  []string{"🥕", "🍆", "🍄"},
		cursor:   0,
		selected: make(map[int]struct{}),
		inputs:   inputs,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "j":
			if m.cursor == len(m.choices)-1 {
				m.cursor = 0
			} else {
				m.cursor++
			}
			return m, nil
		case "k":
			if m.cursor == 0 {
				m.cursor = len(m.choices) - 1
			} else {
				m.cursor--
			}
			return m, nil
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}

	}
	return m, nil
}

func (m model) View() string {
	s := "What do we eat today??\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)

	}
	return s
}

func main() {
	p := tea.NewProgram(initialMode())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
