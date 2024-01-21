package tui

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	configFolder = ".cert-check"
	configFile   = "config.yaml"
)

func InitProject(st State) (tea.Model, tea.Cmd) {
	m := model{}

	if st.ShouldRedail {
		// on startup
		configPath := GetConfigPath(configFolder, configFile)
		if !FileExists(configPath) {
			err := CreateConfig(configFolder, configFile)
			if err != nil {
				log.Fatal(err)
			}
		}

		m.ssls = []ssl{}

		// TODO: merge GetSavedDomains and DialDomains in single function
		savedDomains, err := GetSavedDomains(configPath)
		if err != nil {
			log.Fatal(err)
		}

		DialDomains(&m.ssls, savedDomains)
		if err != nil {
			m.err = err
		}

	} else {
		// coming from entry
		m.ssls = st.ssls
	}

	// populate list
	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(6) // height per item
	m.list = list.New([]list.Item{}, delegate, st.width, st.height)

	for i, s := range m.ssls {
		m.list.InsertItem(i, s)
	}

	// title
	m.list.Styles.Title = styleTitle
	m.list.Title = "🪜 Cert Check"

	return m, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

// interface for list.Item
func (s ssl) FilterValue() string { return s.domain }
func (s ssl) Title() string       { return s.domain }

func (s ssl) Description() string {
	var b strings.Builder
	b.WriteString(RenderDaysLeft(s.daysLeft))
	b.WriteString("Issued On   : " + s.issuedOn + "\n")
	b.WriteString("Expires On  : " + s.expiresOn + "\n")
	b.WriteString("Issuer      : " + s.issuer + "\n")
	b.WriteString("Common Name : " + s.commonName + "\n")
	return b.String()
}
