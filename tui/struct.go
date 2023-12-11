package tui

import "github.com/charmbracelet/bubbles/textinput"

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
