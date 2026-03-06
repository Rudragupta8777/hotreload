package main

import (
	"testing"
)

func TestIsIgnored(t *testing.T) {
	
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"Ignore git folder", ".git/config", true},
		{"Ignore node_modules", "project/node_modules/library", true},
		{"Ignore bin folder", "./bin/server.exe", true},
		{"Allow standard go file", "main.go", false},
		{"Allow nested go file", "internal/api/handler.go", false},
	}

	// Loop through and run each test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isIgnored(tt.path)
			if result != tt.expected {
				t.Errorf("isIgnored(%q) = %v; want %v", tt.path, result, tt.expected)
			}
		})
	}
}