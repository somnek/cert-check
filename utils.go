package main

import "strings"

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
func Sanitize(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ToLower(s)
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
