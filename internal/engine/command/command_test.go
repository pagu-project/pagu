package command

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputBox_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    InputBox
		expected string
		wantErr  bool
	}{
		{"Marshal Text", InputBoxText, `"Text"`, false},
		{"Marshal MultilineText", InputBoxMultilineText, `"MultilineText"`, false},
		{"Marshal Integer", InputBoxInteger, `"Integer"`, false},
		{"Marshal Float", InputBoxFloat, `"Float"`, false},
		{"Marshal File", InputBoxFile, `"File"`, false},
		{"Marshal Toggle", InputBoxToggle, `"Toggle"`, false},
		{"Marshal Choice", InputBoxChoice, `"Choice"`, false},
		{"Marshal Unknown", InputBox(999), ``, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.JSONEq(t, tt.expected, string(data))
			}
		})
	}
}

func TestInputBox_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected InputBox
		wantErr  bool
	}{
		{"Unmarshal Text", `"Text"`, InputBoxText, false},
		{"Unmarshal MultilineText", `"MultilineText"`, InputBoxMultilineText, false},
		{"Unmarshal Integer", `"Integer"`, InputBoxInteger, false},
		{"Unmarshal Float", `"Float"`, InputBoxFloat, false},
		{"Unmarshal File", `"File"`, InputBoxFile, false},
		{"Unmarshal Toggle", `"Toggle"`, InputBoxToggle, false},
		{"Unmarshal Choice", `"Choice"`, InputBoxChoice, false},
		{"Unmarshal Unknown", `"Unknown"`, InputBox(0), true},
		{"Unmarshal Invalid JSON", `123`, InputBox(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var box InputBox
			err := json.Unmarshal([]byte(tt.input), &box)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, box)
			}
		})
	}
}
