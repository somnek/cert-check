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
	configPath := GetConfigPath(configFolder, configFile)
	if !FileExists(configPath) {
		err := CreateConfig(configFolder, configFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	var initSsls []ssl

	// mock data
	initSsls = append(initSsls, ssl{
		domain:     "google.com",
		issuedOn:   "2021-01-01",
		expiresOn:  "2022-01-01",
		issuer:     "Let's Encrypt",
		commonName: "google.com",
	})
	initSsls = append(initSsls, ssl{
		domain:     "facebook.com",
		issuedOn:   "2021-01-01",
		expiresOn:  "2022-01-01",
		issuer:     "Let's Encrypt",
		commonName: "facebook.com",
	})

	// savedDomains, err := GetSavedDomains(configPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = DialDomains(&initSsls, savedDomains)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(5)

	m.list = list.New([]list.Item{}, delegate, st.width, st.height)
	m.list.Title = "🪜 Cert Check"

	for i, s := range initSsls {
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
