package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	configFolder = ".cert-check"
	configFile   = "config.yaml"
)

type model struct {
	input textinput.Model
	ssls  []ssl
	err   error
	logs  string
	page  int
}

type config struct {
	Domains []string `yaml:"domains"`
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
	configPath := getConfigPath(configFolder, configFile)
	if !fileExists(configPath) {
		createConfig(configFolder, configFile)
	}

	var initSsls []ssl
	err := getSavedDomains(&initSsls, configPath)
	if err != nil {
		log.Fatal(err)
	}

	t := textinput.New()
	t.Placeholder = "Type here..."
	t.Focus()
	t.CharLimit = 200
	t.Width = 200

	return model{
		input: t,
		ssls:  initSsls,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// if page is 0, switch to 1, else 0
			m.page = (m.page + 1) % 2
			return m, nil
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			info, err := getInfo(m.input.Value())
			if err != nil {
				m.err = err
				return m, nil
			}

			// prepend to slice
			m.ssls = append([]ssl{info}, m.ssls...)
			m.input.SetValue("")
			m.input.Blur()
			m.page = 0
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
		for _, ssl := range m.ssls {
			s += fmt.Sprintf("%s\n", ssl.domain)
			s += fmt.Sprintf("Issued On   : %s\n", ssl.issuedOn)
			s += fmt.Sprintf("Expires On  : %s\n", ssl.expiresOn)
			s += fmt.Sprintf("Issuer      : %s\n", ssl.issuer)
			s += fmt.Sprintf("Common Name : %s\n\n", ssl.commonName)
		}
	} else {
		s += "Enter a domain name: \n\n"
		s += m.input.View() + "\n\n"
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
