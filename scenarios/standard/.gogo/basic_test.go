package main

import "testing"

// Who cares what we're testing, as long as we have a test file
// This should not show up in GoGo, and it should not choke when test files exist
func TestToUpper(t *testing.T) {
	// Test cases
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "HELLO"},
		{"world", "WORLD"},
	}

	toUpper := func(s string) string {
		result := ""
		for _, char := range s {
			if char >= 'a' && char <= 'z' {
				result += string(char - 32)
			} else {
				result += string(char)
			}
		}
		return result
	}

	for _, test := range tests {
		if got := toUpper(test.input); got != test.expected {
			t.Errorf("toUpper(%s) = %s; want %s", test.input, got, test.expected)
		}
	}
}
