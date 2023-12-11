package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
)

const (
	home page = iota
	form
)

type page int

type model struct {
	cursor int
	input  textinput.Model
	ssls   []ssl
	err    error
	logs   string
	page   page
}

type userConfig struct {
	Domains []string `yaml:"domains"`
}

type ssl struct {
	domain     string
	issuedOn   string
	expiresOn  string
	issuer     string
	commonName string
}
