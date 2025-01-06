package command

import (
	"fmt"
	"html/template"
	"slices"
	"strings"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/color"
)

const failedTemplate = `
**Operation failed**
{{.reason}}
`

const errorTemplate = `
**An error occurred**
{{.err}}
`

var (
	TargetMaskMainnet   = 1
	TargetMaskTestnet   = 2
	TargetMaskModerator = 4

	TargetMaskAll = TargetMaskMainnet | TargetMaskTestnet | TargetMaskModerator
)

type InputBox int

const (
	InputBoxText InputBox = iota
	InputBoxMultilineText
	InputBoxInteger
	InputBoxFloat
	InputBoxFile
	InputBoxToggle
	InputBoxChoice
)

type Choice struct {
	Name  string
	Value int
}

type Args struct {
	Name     string
	Desc     string
	InputBox InputBox
	Optional bool
	Choices  []Choice
}

type HandlerFunc func(caller *entity.User, cmd *Command, args map[string]string) CommandResult

type Command struct {
	Emoji       string
	Color       color.ColorCode
	Name        string
	Help        string
	Args        []Args
	AppIDs      []entity.PlatformID
	SubCommands []*Command
	Middlewares []MiddlewareFunc
	Handler     HandlerFunc
	TargetFlag  int
}

type CommandResult struct {
	Color      color.ColorCode
	Title      string
	Message    string
	Successful bool
}

func (cmd *Command) RenderFailedTemplate(reason string) CommandResult {
	msg, _ := cmd.executeTemplate(failedTemplate, map[string]any{"reason": reason})

	return CommandResult{
		Color:      cmd.Color,
		Title:      fmt.Sprintf("%v %v", cmd.Name, cmd.Emoji),
		Message:    msg,
		Successful: false,
	}
}

func (cmd *Command) RenderErrorTemplate(err error) CommandResult {
	msg, _ := cmd.executeTemplate(errorTemplate, map[string]any{"err": err})

	return CommandResult{
		Color:      cmd.Color,
		Title:      fmt.Sprintf("%v %v", cmd.Name, cmd.Emoji),
		Message:    msg,
		Successful: false,
	}
}

func (cmd *Command) RenderResultTemplate(templateContent string, keyvals ...any) CommandResult {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "!MISSING-VALUE!")
	}

	data := make(map[string]any)
	for i := 0; i < len(keyvals); i += 2 {
		key := keyvals[i].(string)
		val := keyvals[i+1]

		data[key] = val
	}

	msg, err := cmd.executeTemplate(templateContent, data)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	return CommandResult{
		Color:      cmd.Color,
		Title:      fmt.Sprintf("%v %v", cmd.Name, cmd.Emoji),
		Message:    msg,
		Successful: true,
	}
}

func (*Command) executeTemplate(templateContent string, data map[string]any) (string, error) {
	tmpl, _ := template.New("template").Parse(templateContent)

	var bldr strings.Builder
	err := tmpl.Execute(&bldr, data)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(bldr.String()), nil
}

// Deprecated: Use RenderResultTemplate.
func (cmd *Command) SuccessfulResult(msg string) CommandResult {
	return cmd.SuccessfulResultF("%s", msg)
}

// Deprecated: Use RenderResultTemplate.
func (cmd *Command) SuccessfulResultF(msg string, a ...any) CommandResult {
	return CommandResult{
		Color:      cmd.Color,
		Title:      fmt.Sprintf("%v %v", cmd.Name, cmd.Emoji),
		Message:    fmt.Sprintf(msg, a...),
		Successful: true,
	}
}

// Deprecated: Use RenderFailedTemplate.
func (cmd *Command) FailedResult(msg string) CommandResult {
	return cmd.FailedResultF("%s", msg)
}

// Deprecated: Use RenderFailedTemplate.
func (cmd *Command) FailedResultF(msg string, a ...any) CommandResult {
	return CommandResult{
		Color:      cmd.Color,
		Title:      fmt.Sprintf("%v %v", cmd.Name, cmd.Emoji),
		Message:    fmt.Sprintf(msg, a...),
		Successful: false,
	}
}

// Deprecated: Use RenderErrorTemplate.
func (cmd *Command) ErrorResult(err error) CommandResult {
	return cmd.FailedResultF("An error occurred: %v", err.Error())
}

func (cmd *Command) HelpResult() CommandResult {
	return CommandResult{
		Color:      cmd.Color,
		Title:      fmt.Sprintf("%v %v", cmd.Help, cmd.Emoji),
		Message:    cmd.HelpMessage(),
		Successful: false,
	}
}

func (cmd *Command) HasAppID(appID entity.PlatformID) bool {
	return slices.Contains(cmd.AppIDs, appID)
}

func (cmd *Command) HasSubCommand() bool {
	return len(cmd.SubCommands) > 0 && cmd.SubCommands != nil
}

func (cmd *Command) HelpMessage() string {
	help := cmd.Help
	help += "\n\nAvailable commands:\n"
	for _, sc := range cmd.SubCommands {
		help += fmt.Sprintf("- **%-12s**: %s\n", sc.Name, sc.Help)
	}

	return help
}

func (cmd *Command) AddSubCommand(subCmd *Command) {
	if subCmd.HasSubCommand() {
		subCmd.AddHelpSubCommand()
	}

	cmd.SubCommands = append(cmd.SubCommands, subCmd)
}

func (cmd *Command) AddHelpSubCommand() {
	helpCmd := &Command{
		Name:       "help",
		Help:       fmt.Sprintf("Help for %v command", cmd.Name),
		AppIDs:     entity.AllAppIDs(),
		TargetFlag: TargetMaskAll,
		Handler: func(_ *entity.User, _ *Command, _ map[string]string) CommandResult {
			return cmd.SuccessfulResult(cmd.HelpMessage())
		},
	}

	cmd.AddSubCommand(helpCmd)
}
