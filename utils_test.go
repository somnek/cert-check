package main

import "testing"

func TestExtractHost(t *testing.T) {

	t.Run("clean url", func(t *testing.T) {
		url := "google.ca"
		got := ExtractHost(url)
		want := "google.ca"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("url with http", func(t *testing.T) {
		url := "http://google.ca"
		got := ExtractHost(url)
		want := "google.ca"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("url with https", func(t *testing.T) {
		url := "https://google.ca"
		got := ExtractHost(url)
		want := "google.ca"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("url with path", func(t *testing.T) {
		url := "https://google.ca/search"
		got := ExtractHost(url)
		want := "google.ca"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("url with query params", func(t *testing.T) {
		url := "https://google.ca/search?q=github+copilot"
		got := ExtractHost(url)
		want := "google.ca"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

}
