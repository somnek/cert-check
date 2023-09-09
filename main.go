package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	configFolder     = ".cert-check"
	configFile       = "config.yaml"
	MARGIN           = 2
	PADDING          = 1
	COLOR_PINK       = lipgloss.Color("#E95678")
	COLOR_GRAY_1     = lipgloss.Color("#E3E4DB")
	COLOR_GRAY_2     = lipgloss.Color("#CDCDCD")
	COLOR_GRAY_3     = lipgloss.Color("#B8BABA")
	COLOR_GRAY_4     = lipgloss.Color("#626262")
	COLOR_DARK_GREEN = lipgloss.Color("#57886C")
)

var (
	styleSelected = lipgloss.NewStyle().Foreground(COLOR_PINK)
	styleNormal   = lipgloss.NewStyle().Foreground(COLOR_GRAY_1)
	styleTitle    = lipgloss.NewStyle().Foreground(COLOR_GRAY_2).Background(COLOR_DARK_GREEN).Bold(true).Padding(0, 3).MarginTop(1)
	styleHelper1  = lipgloss.NewStyle().Foreground(COLOR_GRAY_4).MarginLeft(1)
	styleHelper2  = lipgloss.NewStyle().Foreground(COLOR_GRAY_3).MarginLeft(1)
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
	configPath := GetConfigPath(configFolder, configFile)
	if !FileExists(configPath) {
		err := CreateConfig(configFolder, configFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	var initSsls []ssl
	err := LoadSavedDomains(&initSsls, configPath)
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

func updateListPage(msg tea.Msg, m model) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "a", "A":
			m.page = 1
			return m, nil

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

		case "q":
			return m, tea.Quit

			// delete domain from config file
		case "d", "D":
			path := GetConfigPath(configFolder, configFile)
			domain := m.ssls[m.cursor].domain
			if err := DeleteFromConfig(domain, path); err != nil {
				m.err = err
				return m, nil
			}
			m.ssls = Delete(m.ssls, m.cursor)
			m.logs += "deleted " + domain + "\n"
			return m, nil

		// save new domain to config file
		case "s", "S":
			path := GetConfigPath(configFolder, configFile)
			domain := m.ssls[m.cursor].domain
			savedDomains, err := GetSavedDomains(path)
			if err != nil {
				m.err = err
			}
			idx := Find(savedDomains, domain)
			if idx == -1 {
				// save domain to config file if it doesn't already exist
				if err := SaveDomain(domain, path); err != nil {
					m.err = err
				}
			} else {
				m.logs += domain + " is already added\n"
			}
			return m, nil
		}
	}
	return m, nil
}

func updateInputPage(msg tea.Msg, m model, cmd *tea.Cmd) (tea.Model, tea.Cmd) {
	m.input, *cmd = m.input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			inputSanitized := Sanitize(m.input.Value())
			info, err := GetInfo(inputSanitized)
			if err != nil {
				m.err = err
				return m, nil
			}

			// check if domain already exists
			domains := ExtractField(m.ssls, "domain")
			idx := Find(domains, inputSanitized)

			if idx != -1 {
				// replace existing domain's info
				m.ssls[idx] = info
			} else {
				// prepend to slice
				m.ssls = append([]ssl{info}, m.ssls...)
			}

			m.input.SetValue("")
			m.input.Blur()
			m.page = 0
			return m, nil
		}
	}
	return m, *cmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		// switch between pages
		// if page is 0, switch to 1, else 0
		case "tab":
			m.page = (m.page + 1) % 2
			if m.page == 0 {
				m.input.Blur()
			} else {
				m.input.SetValue("")
				m.input.Focus()
			}
			m.err = nil
			return m, nil

		case "ctrl+c", "esc":
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	if m.page == 0 {
		return updateListPage(msg, m)
	} else {
		return updateInputPage(msg, m, &cmd)
	}
}

func (m model) View() string {
	var b strings.Builder
	title := "Cert Check"

	b.WriteString(styleTitle.Render(title))
	b.WriteString("\n\n")

	if m.page == 0 {
		// page 0: domain view

		// logs
		if m.logs != "" {
			b.WriteString(m.logs)
			b.WriteString("\n")
		}

		// domain list
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

			b.WriteString(style.Render(fmt.Sprintf("%s %s", cursor, ssl.domain)))
			b.WriteString("\n")
			b.WriteString(style.Render(fmt.Sprintf("%s Issued On   : %s", cursor, ssl.issuedOn)))
			b.WriteString("\n")
			b.WriteString(style.Render(fmt.Sprintf("%s Expires On  : %s", cursor, ssl.expiresOn)))
			b.WriteString("\n")
			b.WriteString(style.Render(fmt.Sprintf("%s Issuer      : %s", cursor, ssl.issuer)))
			b.WriteString("\n")
			b.WriteString(style.Render(fmt.Sprintf("%s Common Name : %s", cursor, ssl.commonName)))
			b.WriteString("\n\n")
		}

		// print errors
		if m.err != nil {
			b.WriteString(m.err.Error())
		}

		renderPageControls(&b, m.page)

	} else {
		// page 1: input view
		b.WriteString("Enter a domain name: \n\n")
		b.WriteString(m.input.View() + "\n\n")

		// print errors
		if m.err != nil {
			b.WriteString(m.err.Error() + "\n")
		}

		renderPageControls(&b, m.page)

	}

	return b.String()
}

func main() {
	p := tea.NewProgram(initialMode(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
