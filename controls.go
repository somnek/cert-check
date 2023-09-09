package main

import "strings"

// TODO: abstract this into a struct
func renderPageControls(b *strings.Builder, page int) {
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
	} else if page == 1 {
		b.WriteString(styleHelper1.Render("ctrl+c: "))
		b.WriteString(styleHelper2.Render("quit "))
		b.WriteString(styleHelper2.Render("• "))
		b.WriteString(styleHelper1.Render("tab: "))
		b.WriteString(styleHelper2.Render("switch page\n"))
	}
}
