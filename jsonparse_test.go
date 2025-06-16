package main

import "testing"

func TestParseJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{" false ", true},
		{"null", true},
		{"invalid", false},
		{"", false},
		{"tru", false},
		{"falsee", false},
		{"nullish", false},
		{"\n\ttrue\r", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseJSON(tt.input)
			if result != tt.expected {
				t.Errorf("parseJSON(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
