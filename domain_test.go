package main

import (
	"testing"
)

func TestGetInfo(t *testing.T) {
	_, err := getInfo("google.ca")
	if err != nil {
		t.Errorf("getInfo() failed")
	}
}
