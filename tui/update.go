package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

func updateListPage(msg tea.Msg, m model) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "a", "A":
			m.input.SetValue("")
			m.input.Focus()
			m.page = 1
			return m, nil

		case "j":
			if m.cursor < len(m.ssls)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
			m.logs = ""
			return m, nil

		case "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.ssls) - 1
			}
			m.logs = ""
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
			m.logs = "deleted " + domain + "\n"
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

				// save domain to config file
				path := GetConfigPath(configFolder, configFile)
				if err := SaveDomain(inputSanitized, path); err != nil {
					m.err = err
				}
				m.logs = "added " + inputSanitized + "\n"
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
