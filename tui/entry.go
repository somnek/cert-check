package tui

import (
	"errors"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type entry struct {
	input  textinput.Model
	ssls   []ssl
	width  int
	height int
}

func (m entry) Init() tea.Cmd {
	return nil
}

func InitEntry(ssls []ssl) *entry {
	// text input
	t := textinput.New()
	t.Placeholder = "Type here..."
	t.Focus()
	t.CharLimit = 200
	t.Width = 200

	m := entry{
		input: t,
	}
	m.ssls = ssls
	return &m
}

func (m entry) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.input, cmd = m.input.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Tab):
			st := State{width: m.width, height: m.height, ssls: m.ssls, ShouldRedail: false}
			return InitProject(st)
		case key.Matches(msg, keys.Enter):
			m.input.Validate = func(s string) error {
				if s == " " {
					return errors.New("cannot be empty")
				}
				return nil
			}

			// dail
			domain := Sanitize(m.input.Value())
			info, err := GetInfo(domain)
			if err != nil {
				log.Fatal(err)
			}

			m.ssls = append(m.ssls, info)

			// Save
			// check if domain already exists
			domains := ExtractField(m.ssls, "domain")
			idx := Find(domains, info.domain)

			if idx != -1 {
				// replace existing domain's info
				m.ssls[idx] = info
			} else {
				// prepend to slice
				m.ssls = append([]ssl{info}, m.ssls...)

				// save domain to config file
				path := GetConfigPath(configFolder, configFile)
				if err := SaveDomain(info.domain, path); err != nil {
					log.Fatal(err)
				}

			}

			st := State{
				width:        m.width,
				height:       m.height,
				newSsl:       info,
				ssls:         m.ssls,
				ShouldRedail: false,
			}
			return InitProject(st)
		}
	}
	return m, cmd
}

func (m entry) View() string {
	return m.input.View()
}
