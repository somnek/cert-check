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

func InitProject(redail bool) (tea.Model, tea.Cmd) {
	m := model{}
	configPath := GetConfigPath(configFolder, configFile)
	if !FileExists(configPath) {
		err := CreateConfig(configFolder, configFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	var initSsls []ssl
	savedDomains, err := GetSavedDomains(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if redail {
		err = DialDomains(&initSsls, savedDomains)
		if err != nil {
			log.Fatal(err)
		}
	}

	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(4)

	m.list = list.New([]list.Item{}, delegate, 8, 8)
	m.list.Title = "🪜 Cert Check"

	for i, s := range m.ssls {
		m.list.InsertItem(i, s)
	}
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
	b.WriteString("Issued On   : " + s.issuedOn + "\n")
	b.WriteString("Expires On  : " + s.expiresOn + "\n")
	b.WriteString("Issuer      : " + s.issuer + "\n")
	b.WriteString("Common Name : " + s.commonName + "\n")
	return b.String()
}
