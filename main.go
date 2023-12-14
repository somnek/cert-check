package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/somnek/cert-check/tui"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	st := tui.State{}
	m, _ := tui.InitProject(st)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
