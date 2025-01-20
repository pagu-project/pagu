package utils_test

import (
	"testing"

	"github.com/pagu-project/pagu/pkg/utils"
)

func TestCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Single word lowercase", "word", "word"},
		{"Single word uppercase", "WORD", "word"},
		{"Hyphen-separated", "hello-world", "helloWorld"},
		{"Underscore-separated", "hello_world", "helloWorld"},
		{"Mixed delimiters", "hello-world_example", "helloWorldExample"},
		{"Multiple spaces", "hello   world", "helloWorld"},
		{"Leading and trailing spaces", "  hello world  ", "helloWorld"},
		{"Mixed case input", "HeLLo-WoRLD", "helloWorld"},
		{"Numbers in input", "hello-world-123", "helloWorld123"},
		{"Special characters", "hello-world!", "helloWorld"},
		{"Consecutive delimiters", "hello--world__example", "helloWorldExample"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.CamelCase(test.input)
			if result != test.expected {
				t.Errorf("CamelCase(%q) = %q; expected %q", test.input, result, test.expected)
			}
		})
	}
}
