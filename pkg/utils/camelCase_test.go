package utils_test

import (
	"testing"

	"github.com/pagu-project/pagu/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Single word", "hello", "hello"},
		{"Multiple words", "hello world", "helloWorld"},
		{"Hyphenated words", "hello-world", "helloWorld"},
		{"Mixed case", "Hello-World test", "helloWorldTest"},
		{"Leading spaces", "   hello world", "helloWorld"},
		{"Trailing spaces", "hello world   ", "helloWorld"},
		{"Extra spaces between words", "hello   world", "helloWorld"},
		{"Empty input", "", ""},
		{"Only spaces", "   ", ""},
		{"Hyphenated with spaces", "hello - world", "helloWorld"},
		{"Special characters", "hello-world123", "helloWorld123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.CamelCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
