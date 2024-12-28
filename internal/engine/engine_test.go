package engine

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantCmds []string
		wantArgs map[string]string
		wantErr  error
	}{
		{
			name:     "valid input with commands and arguments",
			input:    "command1 command2 --arg1=val1 --arg2=val2",
			wantCmds: []string{"command1", "command2"},
			wantArgs: map[string]string{"arg1": "val1", "arg2": "val2"},
			wantErr:  nil,
		},
		{
			name:     "input with no arguments",
			input:    "command1 command2",
			wantCmds: []string{"command1", "command2"},
			wantArgs: map[string]string{},
			wantErr:  nil,
		},
		{
			name:     "input with arguments only",
			input:    "--arg1=val1 --arg2=val2",
			wantCmds: []string{},
			wantArgs: map[string]string{"arg1": "val1", "arg2": "val2"},
			wantErr:  nil,
		},
		{
			name:     "invalid argument format (missing =)",
			input:    "command1 --arg1",
			wantCmds: nil,
			wantArgs: nil,
			wantErr:  fmt.Errorf("invalid argument format: --arg1"),
		},
		{
			name:     "invalid argument format (empty key)",
			input:    "command1 --=val1",
			wantCmds: nil,
			wantArgs: nil,
			wantErr:  fmt.Errorf("invalid argument format: --=val1"),
		},
		{
			name:     "invalid argument format (empty value)",
			input:    "command1 --arg1=",
			wantCmds: nil,
			wantArgs: nil,
			wantErr:  fmt.Errorf("invalid argument format: --arg1="),
		},
		{
			name:     "empty input",
			input:    "",
			wantCmds: nil,
			wantArgs: nil,
			wantErr:  errors.New("input string cannot be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmds, gotArgs, gotErr := parseCommand(tt.input)

			// Compare commands
			assert.Equal(t, tt.wantCmds, gotCmds, "commands mismatch: got %v, want %v", gotCmds, tt.wantCmds)

			// Compare arguments
			assert.Equal(t, tt.wantArgs, gotArgs, "arguments mismatch: got %v, want %v", gotArgs, tt.wantArgs)

			// Compare errors
			if tt.wantErr != nil {
				assert.ErrorContains(t, gotErr, tt.wantErr.Error(), "error mismatch: got %v, want %v", gotErr, tt.wantErr)
			} else {
				assert.NoError(t, gotErr, "error mismatch: got %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
