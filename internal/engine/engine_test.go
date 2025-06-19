package engine

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/engine/command/phoenix"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/testsuite"
	"github.com/pagu-project/pagu/pkg/client"
	"github.com/pagu-project/pagu/pkg/mailer"
	"github.com/pagu-project/pagu/pkg/nowpayments"
	"github.com/pagu-project/pagu/pkg/wallet"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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
			name:     "repeated args",
			input:    "command1 command2 --arg1=val1 --arg1=val1",
			wantCmds: []string{"command1", "command2"},
			wantArgs: map[string]string{"arg1": "val1"},
			wantErr:  nil,
		},
		{
			name:     "arguments with quotation",
			input:    "command1 --arg1='val1' --arg2=\"val2\"",
			wantCmds: []string{"command1"},
			wantArgs: map[string]string{"arg1": "val1", "arg2": "val2"},
			wantErr:  nil,
		},
		{
			name:     "arguments with quotation inside value",
			input:    "command1 --arg1='val ' 1' --arg2=\"val \" 2\"",
			wantCmds: []string{"command1"},
			wantArgs: map[string]string{"arg1": "val ' 1", "arg2": "val \" 2"},
			wantErr:  nil,
		},
		{
			name:     "arguments with quotation and spaces",
			input:    "command1 --arg1='val 1' --arg2=\"val 2\" --arg3=val3",
			wantCmds: []string{"command1"},
			wantArgs: map[string]string{"arg1": "val 1", "arg2": "val 2", "arg3": "val3"},
			wantErr:  nil,
		},
		{
			name:     "arguments with = inside value",
			input:    "command1 --arg1='val=1'",
			wantCmds: []string{"command1"},
			wantArgs: map[string]string{"arg1": "val=1"},
			wantErr:  nil,
		},
		{
			name:     "extra spaces",
			input:    "command1    command2   --arg1=' val 1'    --arg2=\"val 2 \"",
			wantCmds: []string{"command1", "command2"},
			wantArgs: map[string]string{"arg1": "val 1", "arg2": "val 2"},
			wantErr:  nil,
		},
		{
			name:     "with tabs",
			input:    "command1		--arg1='val	1	'	--arg2=\"	val	2\"",
			wantCmds: []string{"command1"},
			wantArgs: map[string]string{"arg1": "val 1", "arg2": "val 2"},
			wantErr:  nil,
		},
		{
			name:     "argument with empty value",
			input:    "command1 --arg1=\"\"",
			wantCmds: []string{"command1"},
			wantArgs: map[string]string{"arg1": ""},
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
			wantCmds: []string{"command1"},
			wantArgs: map[string]string{"arg1": "true"},
			wantErr:  nil,
		},
		{
			name:     "invalid argument format (empty key)",
			input:    "command1 --=val1",
			wantCmds: nil,
			wantArgs: nil,
			wantErr:  errors.New("invalid argument format: --=val1"),
		},
		{
			name:     "invalid argument format (empty value)",
			input:    "command1 --arg1=",
			wantCmds: []string{"command1"},
			wantArgs: map[string]string{"arg1": ""},
			wantErr:  nil,
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
			gotCmds, gotArgs, gotErr := parseInput(tt.input)

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

func isLowerCase(s string) bool {
	return s == strings.ToLower(s)
}

func endsWithPeriod(s string) bool {
	return strings.HasSuffix(s, ".")
}

func TestCheckCommandsAndArgs(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	ctrl := gomock.NewController(t)

	ctx := context.Background()
	testDB := ts.MakeTestDB()
	mockClientManager := client.NewMockIManager(ctrl)
	mockWallet := wallet.NewMockIWallet(ctrl)
	mockMailer := mailer.NewMockIMailer(ctrl)
	mockNowPayments := nowpayments.NewMockINowPayments(ctrl)
	cfg := &Config{
		Phoenix: phoenix.Config{
			PrivateKey: "TSECRET1RZSMS2JGNFLRU26NHNQK3JYTD4KGKLGW4S7SG75CZ057SR7CE8HUSG5MS3Z",
		},
	}
	eng := newBotEngine(ctx, entity.BotID_CLI, cfg,
		testDB, mockClientManager, mockWallet, mockMailer, mockNowPayments)

	var checkCommands func(cmds []*command.Command)

	checkCommands = func(cmds []*command.Command) {
		for _, cmd := range cmds {
			if cmd.Help == "" {
				t.Errorf("Command has no help: %s", cmd.Name)
			}

			if !isLowerCase(cmd.Name) {
				t.Errorf("Command name is not lowercase: %s", cmd.Name)
			}

			if endsWithPeriod(cmd.Help) {
				t.Errorf("Command help should not end with a period: %s", cmd.Help)
			}

			for _, arg := range cmd.Args {
				if !isLowerCase(arg.Name) {
					t.Errorf("Argument name is not lowercase: %s", cmd.Name)
				}

				if endsWithPeriod(arg.Desc) {
					t.Errorf("Argument desc should not end with a period: %s", cmd.Help)
				}
			}

			checkCommands(cmd.SubCommands)
		}
	}

	checkCommands(eng.Commands())
}
