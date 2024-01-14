package tui

import (
	"github.com/charmbracelet/bubbles/list"
)

type model struct {
	list list.Model
	ssls []ssl
	err  error
	logs string
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

// channel response from dialing domain
type chDailRes struct {
	ssl ssl
	err error
}
