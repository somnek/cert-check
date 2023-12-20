package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/somnek/cert-check/tui"
)

func main() {
	st := tui.State{ShouldRedail: true}
	m, _ := tui.InitProject(st)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
