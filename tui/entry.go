package tui

import (
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

	// remove unnecessary keys for entry page
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
		case key.Matches(msg, keys.Tab, keys.Esc):
			st := State{width: m.width, height: m.height, ssls: m.ssls, ShouldRedail: false}
			return InitProject(st)
		case key.Matches(msg, keys.Enter):
			// dail
			domain := Sanitize(m.input.Value())
			ch := make(chan chDailRes)
			go GetInfo(domain, ch)
			res := <-ch
			if res.err != nil {
				m.err = res.err
				m.input.SetValue("")
				m.input.Focus()
				return m, nil
			}

			info := res.ssl

			// Save
			// check if domain already exists
			domains := ExtractField(m.ssls, "domain")
			idx := Find(domains, info.domain)

			if idx != -1 {
				// exists: move to top of slice
				m.ssls = Delete(m.ssls, idx)
				m.ssls = append([]ssl{info}, m.ssls...)
			} else {
				// new: prepend to slice
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

	// title
	title := "🪜 Cert Check"
	sb.WriteString(styleTitle.Render(title))
	sb.WriteString("\n\n")

	// list
	sb.WriteString(m.input.View())

	if m.err != nil {
		sb.WriteString(fmt.Sprintf("\n\n🚫 erorr: %v\n", m.err.Error()))
	}

	// help
	help := "enter submit • esc back"
	height := 8 - strings.Count(sb.String(), "\n") - strings.Count(help, "\n")
	sb.WriteString(strings.Repeat("\n", height))
	sb.WriteString(styleHelper1.Render(help))

	return stylePadEntry.Render(sb.String())
}
