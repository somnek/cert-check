package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	configFolder = ".cert-check"
	configFile   = "config.yaml"
)

var (
	styleSelected = lipgloss.NewStyle().Foreground(lipgloss.Color("#E95678"))
	styleNormal   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	styleTitle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA")).Bold(true).Padding(0, 3).MarginTop(1)
)

type model struct {
	cursor int
	input  textinput.Model
	ssls   []ssl
	err    error
	logs   string
	page   int
}

type userConfig struct {
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

		case "j":
			if m.cursor < len(m.ssls)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
			return m, nil

		case "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.ssls) - 1
			}
			return m, nil

		case "tab":
			// if page is 0, switch to 1, else 0
			m.page = (m.page + 1) % 2
			if m.page == 1 {
				m.input.Focus()
			} else {
				m.input.Blur()
			}
			return m, nil

		case "ctrl+c", "q":
			return m, tea.Quit

		// save new domain to config file
		case "A":
			path := getConfigPath(configFolder, configFile)
			domain := m.ssls[m.cursor].domain
			err := saveDomain(domain, path)
			if err != nil {
				m.err = err
				return m, nil
			}

		case "enter":

			// does nothing on page 0
			if m.page == 0 {
				return m, nil
			} else {
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
	title := "Cert Check"
	s += styleTitle.Render(title)
	s += "\n\n"

	if m.page == 0 {
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

			s += style.Render(fmt.Sprintf("%s %s", cursor, ssl.domain))
			s += "\n"
			s += style.Render(fmt.Sprintf("%s Issued On   : %s", cursor, ssl.issuedOn))
			s += "\n"
			s += style.Render(fmt.Sprintf("%s Expires On  : %s", cursor, ssl.expiresOn))
			s += "\n"
			s += style.Render(fmt.Sprintf("%s Issuer      : %s", cursor, ssl.issuer))
			s += "\n"
			s += style.Render(fmt.Sprintf("%s Common Name : %s", cursor, ssl.commonName))
			s += "\n\n"
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
	p := tea.NewProgram(initialMode())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
