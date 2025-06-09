// Package command provides structures for handling commands.
package command

import (
	"fmt"
	"html/template"
	"slices"
	"strings"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/version"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/pagu-project/pagu/pkg/utils"
)

const failedTemplate = `
**Operation failed:**
{{.reason}}
`

const errorTemplate = `
**An error occurred:**
{{.err}}
`

type InputBox int

const (
	InputBoxText          InputBox = 1
	InputBoxMultilineText InputBox = 2
	InputBoxInteger       InputBox = 3
	InputBoxFloat         InputBox = 4
	InputBoxFile          InputBox = 5
	InputBoxToggle        InputBox = 6
	InputBoxChoice        InputBox = 7
)

var inputBoxToString = map[InputBox]string{
	InputBoxText:          "Text",
	InputBoxMultilineText: "MultilineText",
	InputBoxInteger:       "Integer",
	InputBoxFloat:         "Float",
	InputBoxFile:          "File",
	InputBoxToggle:        "Toggle",
	InputBoxChoice:        "Choice",
}

func (ib InputBox) String() string {
	str, ok := inputBoxToString[ib]
	if ok {
		return str
	}

	return fmt.Sprintf("%d", ib)
}

func (ib InputBox) MarshalYAML() (any, error) {
	return utils.MarshalEnum(ib, inputBoxToString)
}

func (ib *InputBox) UnmarshalYAML(unmarshal func(any) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	val, err := utils.UnmarshalEnum(str, inputBoxToString)
	if err != nil {
		return err
	}
	*ib = val

	return nil
}

type Choice struct {
	Name  string `yaml:"name"`
	Desc  string `yaml:"desc"`
	Value string `yaml:"value"`
}

type Args struct {
	Name     string   `yaml:"name"`
	Desc     string   `yaml:"desc"`
	InputBox InputBox `yaml:"input_box"`
	Optional bool     `yaml:"optional"`
	Choices  []Choice `yaml:"choices"`
}

type HandlerFunc func(caller *entity.User, cmd *Command, args map[string]string) CommandResult

type Command struct {
	Name           string         `yaml:"name"`
	Emoji          string         `yaml:"emoji"`
	Active         bool           `yaml:"active"`
	Help           string         `yaml:"help"`
	Args           []*Args        `yaml:"args"`
	SubCommands    []*Command     `yaml:"sub_commands"`
	ResultTemplate string         `yaml:"result_template"`
	TargetBotIDs   []entity.BotID `yaml:"target_bot_ids"`
	Handler        HandlerFunc    `yaml:"-"`
}

type CommandResult struct {
	Title      string
	Message    string
	Successful bool
}

func (cmd *Command) RenderInternalFailure() CommandResult {
	return cmd.RenderFailedTemplate("An internal error happened. Please try later")
}

func (cmd *Command) RenderFailedTemplateF(reason string, a ...any) CommandResult {
	return cmd.RenderFailedTemplate(fmt.Sprintf(reason, a...))
}

func (cmd *Command) RenderFailedTemplate(reason string) CommandResult {
	msg := cmd.executeTemplate(failedTemplate, map[string]any{"reason": reason})

	return CommandResult{
		Title:      fmt.Sprintf("%v %v", cmd.Name, cmd.Emoji),
		Message:    msg,
		Successful: false,
	}
}

func (cmd *Command) RenderErrorTemplate(err error) CommandResult {
	msg := cmd.executeTemplate(errorTemplate, map[string]any{"err": err})

	return CommandResult{
		Title:      fmt.Sprintf("%v %v", cmd.Name, cmd.Emoji),
		Message:    msg,
		Successful: false,
	}
}

func (cmd *Command) RenderResultTemplate(keyvals ...any) CommandResult {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "!MISSING-VALUE!")
	}

	data := make(map[string]any)
	for i := 0; i < len(keyvals); i += 2 {
		key := keyvals[i].(string)
		val := keyvals[i+1]

		data[key] = val
	}

	msg := cmd.executeTemplate(cmd.ResultTemplate, data)

	return CommandResult{
		Title:      fmt.Sprintf("%v %v", cmd.Name, cmd.Emoji),
		Message:    msg,
		Successful: true,
	}
}

func (*Command) executeTemplate(templateContent string, data map[string]any) string {
	funcMap := template.FuncMap{
		"fixed": func(width int, s string) string {
			if len(s) > width {
				return s[:width]
			}

			return s + strings.Repeat(" ", width-len(s))
		},
	}

	templateContent = strings.ReplaceAll(templateContent, "”", "`")
	tmpl, err := template.New("template").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		log.Error("unable to parse template", "error", err)
	}

	var bldr strings.Builder
	err = tmpl.Execute(&bldr, data)
	if err != nil {
		log.Error("unable to parse template", "error", err)
	}

	return strings.TrimSpace(bldr.String())
}

// Deprecated: Use RenderResultTemplate.
func (cmd *Command) SuccessfulResult(msg string) CommandResult {
	return cmd.SuccessfulResultF("%s", msg)
}

// Deprecated: Use RenderResultTemplate.
func (cmd *Command) SuccessfulResultF(msg string, a ...any) CommandResult {
	return CommandResult{
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
		Title:      fmt.Sprintf("%v %v", cmd.Name, cmd.Emoji),
		Message:    fmt.Sprintf(msg, a...),
		Successful: false,
	}
}

// Deprecated: Use RenderErrorTemplate.
func (cmd *Command) ErrorResult(err error) CommandResult {
	return cmd.FailedResultF("An error occurred: %v", err.Error())
}

func (cmd *Command) RenderHelpTemplate() CommandResult {
	const helpCommandTemplate = `
{{.cmd.Help}}

**Usage:**
   {{.cmd.Name}} [subcommand]

**Available Subcommands:**
{{- range .cmd.SubCommands }}
- ”{{.Name }}” {{.Emoji}} {{.Help}}
{{- end}}

Use "{{.cmd.Name}} help --subcommand=[subcommand]" for more information about a subcommand.
`
	msg := cmd.executeTemplate(helpCommandTemplate, map[string]any{"cmd": cmd})

	return CommandResult{
		Title:      fmt.Sprintf("%v %v", cmd.Name, cmd.Emoji),
		Message:    msg,
		Successful: true,
	}
}

func (cmd *Command) NameWithEmoji() string {
	if cmd.Emoji != "" {
		return fmt.Sprintf("%s %s", cmd.Emoji, cmd.Name)
	}

	return cmd.Name
}

func (cmd *Command) CanBeHandledByBot(botID entity.BotID) bool {
	return slices.Contains(cmd.TargetBotIDs, botID)
}

func (cmd *Command) HasSubCommand() bool {
	return len(cmd.SubCommands) > 0 && cmd.SubCommands != nil
}

func (cmd *Command) AddSubCommand(subCmd *Command) {
	if subCmd == nil {
		return
	}

	if subCmd.HasSubCommand() {
		subCmd.AddHelpSubCommand()
	}

	cmd.SubCommands = append(cmd.SubCommands, subCmd)
}

func (cmd *Command) AddHelpSubCommand() {
	helpCmd := &Command{
		Emoji:        "❓",
		Name:         "help",
		Help:         fmt.Sprintf("Help for %v command", cmd.Name),
		TargetBotIDs: entity.AllBotIDs(),
		Handler: func(_ *entity.User, _ *Command, _ map[string]string) CommandResult {
			return cmd.RenderHelpTemplate()
		},
	}

	cmd.AddSubCommand(helpCmd)
}

func (cmd *Command) AddAboutSubCommand() {
	//nolint:dupword // Pagu is duplicated in the about command
	const aboutTemplate = `
## About Pagu

Pagu is a multi-platform bot designed to interact with and monitor the Pactus Blockchain.
It offers real-time network monitoring, block reward estimation, PAC coin market prices,
Phoenix Testnet integration, and More...

🏷️ Version: {{.version}}
🌐 Homepage: https://pagu.bot

---------

## About Pactus

🌐 Website:  https://pactus.org
🔐 Wallet:   https://wallet.pactus.org
🔍 Explorer: https://pacviewer.com
`

	cmd.ResultTemplate = aboutTemplate
	aboutCmd := &Command{
		Emoji:          "📝",
		Name:           "about",
		Help:           "About Pagu",
		TargetBotIDs:   entity.AllBotIDs(),
		ResultTemplate: aboutTemplate,
		Handler: func(_ *entity.User, _ *Command, _ map[string]string) CommandResult {
			return cmd.RenderResultTemplate("version", version.StringVersion())
		},
	}

	cmd.AddSubCommand(aboutCmd)
}
