package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func writeDomainList(b *strings.Builder, s ssl, style lipgloss.Style, cursor string) {
	b.WriteString(style.Render(fmt.Sprintf("%s %s", cursor, s.domain)))
	b.WriteString("\n")
	b.WriteString(style.Render(fmt.Sprintf("%s Issued On   : %s", cursor, s.issuedOn)))
	b.WriteString("\n")
	b.WriteString(style.Render(fmt.Sprintf("%s Expires On  : %s", cursor, s.expiresOn)))
	b.WriteString("\n")
	b.WriteString(style.Render(fmt.Sprintf("%s Issuer      : %s", cursor, s.issuer)))
	b.WriteString("\n")
	b.WriteString(style.Render(fmt.Sprintf("%s Common Name : %s", cursor, s.commonName)))
	b.WriteString("\n\n")
}

// TODO: abstract this into a struct
func writePageControls(b *strings.Builder, page int) {
	if page == 0 {
		b.WriteString(styleHelper1.Render("j/k ↑/↓ "))
		b.WriteString(styleHelper2.Render("select "))
		b.WriteString(styleHelper2.Render("• "))
		b.WriteString(styleHelper1.Render("q: "))
		b.WriteString(styleHelper2.Render("quit "))
		b.WriteString(styleHelper2.Render("• "))
		b.WriteString(styleHelper1.Render("tab: "))
		b.WriteString(styleHelper2.Render("switch page "))
		b.WriteString(styleHelper2.Render("• "))
		b.WriteString(styleHelper1.Render("d: "))
		b.WriteString(styleHelper2.Render("delete "))
		b.WriteString(styleHelper2.Render("• "))
		b.WriteString(styleHelper1.Render("a: "))
		b.WriteString(styleHelper2.Render("add\n"))
	} else if page == 1 {
		b.WriteString(styleHelper1.Render("ctrl+c: "))
		b.WriteString(styleHelper2.Render("quit "))
		b.WriteString(styleHelper2.Render("• "))
		b.WriteString(styleHelper1.Render("tab: "))
		b.WriteString(styleHelper2.Render("switch page\n"))
	}
}
