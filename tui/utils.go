package tui

import (
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

// CalculateDaysLeft returns the number of days left until expiry
func CalculateDaysLeft(expiry time.Time) int {
	return int(time.Until(expiry).Hours() / 24)
	// return int(expiry.Sub(time.Now()).Hours() / 24)
}

// hasScheme returns true if the url has a scheme
func hasScheme(rawUrl string) bool {
	return strings.HasPrefix(rawUrl, "http://") || strings.HasPrefix(rawUrl, "https://")
}

// ExtractHost returns the host from a given url
func ExtractHost(rawUrl string) string {
	if hasScheme(rawUrl) {
		u, err := url.Parse(rawUrl)
		if err != nil {
			log.Fatal(err)
		}
		return u.Host
	}
	return rawUrl
}

// ExtractField returns a slice of values for a given field
func ExtractField(ssls []ssl, f string) []string {
	vals := make([]string, len(ssls))

	for i, ssl := range ssls {
		switch f {
		case "domain":
			vals[i] = ssl.domain
		case "issuedOn":
			vals[i] = ssl.issuedOn
		case "expiresOn":
			vals[i] = ssl.expiresOn
		case "issuer":
			vals[i] = ssl.issuer
		case "commonName":
			vals[i] = ssl.commonName
		}
	}
	return vals
}

// Sanitize removes spaces and converts to lowercase
func Sanitize(raw string) string {
	s := raw
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ToLower(s)
	s = ExtractHost(s)
	return s
}

// Find returns the index of the first instance of value in slice
// or -1 if value is not present in slice
func Find[T comparable](l []T, v T) int {
	for i, e := range l {
		if e == v {
			return i
		}
	}
	return -1
}

// Delete removes the element at index i from slice l
// returns the modified slice
func Delete[T any](l []T, i int) []T {
	if len(l) < 1 {
		return l
	}
	return append(l[:i], l[i+1:]...)
}

// logToFile prints output to debug.log
func logToFile(s ...any) {
	f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(nil)
	}
	defer f.Close()
	log.SetOutput(f)
	for _, t := range s {
		log.Print(t)
	}
}
