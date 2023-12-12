package tui

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	configFolder = ".cert-check"
	configFile   = "config.yaml"
)

func New() *model {
	// text input
	t := textinput.New()
	t.Placeholder = "Type here..."
	t.Focus()
	t.CharLimit = 200
	t.Width = 200

	return &model{
		input: t,
	}

}

func (m *model) initList(width, height int) {
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

	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(4)
	m.list = list.New([]list.Item{}, delegate, width, height)
	m.list.Title = "🪜 Cert Check"
	for i, s := range initSsls {
		m.list.InsertItem(i, s)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (s ssl) FilterValue() string {
	return s.domain
}
func (s ssl) Title() string {
	return s.domain
}
func (s ssl) Description() string {
	var b strings.Builder
	b.WriteString("Issued On   : " + s.issuedOn + "\n")
	b.WriteString("Expires On  : " + s.expiresOn + "\n")
	b.WriteString("Issuer      : " + s.issuer + "\n")
	b.WriteString("Common Name : " + s.commonName + "\n")
	return b.String()
}
