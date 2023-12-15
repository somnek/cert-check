package tui

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type entry struct {
	input  textinput.Model
	ssls   []ssl
	width  int
	height int
	err    error
}

func (m entry) Init() tea.Cmd {
	return nil
}

func InitEntry(ssls []ssl) *entry {
	// text input
	t := textinput.New()
	t.Placeholder = "Enter new domain..."
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
				m.err = err
				m.input.SetValue("")
				m.input.Focus()
				return m, nil
			}

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
	var sb strings.Builder

	sb.WriteString(m.input.View())

	if m.err != nil {
		sb.WriteString(fmt.Sprintf("\n\n🚫 erorr: %v\n", m.err.Error()))
	}

	return sb.String()
}
