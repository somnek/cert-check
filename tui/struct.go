package tui

import (
	"github.com/charmbracelet/bubbles/list"
)

const (
	home page = iota
	form
)

type page int
type errMsg error

type model struct {
	list   list.Model
	cursor int
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

// state data shared that are passed back to initProject
type State struct {
	width        int
	height       int
	ssls         []ssl
	newSsl       ssl
	ShouldRedail bool
}
