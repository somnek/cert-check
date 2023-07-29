package main

import (
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	input textinput.Model
	ssls  []ssl
	err   error
	page  int
	logs  string
}

type errMsg error

type ssl struct {
	domain     string
	issuedOn   string
	expiresOn  string
	org        string
	issuer     string
	commonName string
}

func initialMode() model {
	t := textinput.New()
	t.Placeholder = "Type here"
	t.Focus()
	t.CharLimit = 200
	t.Width = 200

	return model{
		input: t,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			ssl, err := getInfo(m.input.Value())
			if err != nil {
				log.Print("<--\n")
				log.Fatal(err)
				log.Print("-->\n")
				return m, tea.Quit
			}

			m.ssls = append(m.ssls, ssl)
			m.input.SetValue("")
			m.input.Blur()
			m.page = 1
			return m, nil
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var s string
	if m.page == 0 {
		s += "Enter a domain name: \n\n"
		s += m.input.View() + "\n\n"
	} else {
		for _, ssl := range m.ssls {
			s += "bread 🥖"
			s += ssl.issuedOn + "\n"
			s += ssl.expiresOn + "\n"
			s += ssl.issuer + "\n"
			s += ssl.commonName
		}
	}
	return s
}

func main() {
	p := tea.NewProgram(initialMode(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
