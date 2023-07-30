package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	input textinput.Model
	ssls  []ssl
	err   error
	logs  string
}

type errMsg error

type ssl struct {
	domain     string
	issuedOn   string
	expiresOn  string
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
				m.err = err
				return m, nil
			}

			m.ssls = append(m.ssls, ssl)
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
	s += "Enter a domain name: \n\n"
	s += m.input.View() + "\n\n"
	for _, ssl := range m.ssls {
		s += fmt.Sprintf("Domain      : %s\n", ssl.domain)
		s += fmt.Sprintf("Issued On   : %s\n", ssl.issuedOn)
		s += fmt.Sprintf("Expires On  : %s\n", ssl.expiresOn)
		s += fmt.Sprintf("Issuer      : %s\n", ssl.issuer)
		s += fmt.Sprintf("Common Name : %s\n\n", ssl.commonName)
	}

	// print errors
	if m.err != nil {
		s += m.err.Error()
	}

	return s
}

func main() {
	p := tea.NewProgram(initialMode(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
