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
	Emoji          string              `yaml:"emoji"`
	Name           string              `yaml:"name"`
	Help           string              `yaml:"help"`
	Args           []*Args             `yaml:"args"`
	SubCommands    []*Command          `yaml:"sub_commands"`
	ResultTemplate string              `yaml:"result_template"`
	Middlewares    []MiddlewareFunc    `yaml:"-"`
	Handler        HandlerFunc         `yaml:"-"`
	AppIDs         []entity.PlatformID `yaml:"-"`
	TargetFlag     int                 `yaml:"-"`
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

**Available subcommands:**
   {{- range .cmd.SubCommands }}
   <pre>{{.Name | fixed 15 }}</pre> {{.Emoji}} {{.Help}}
   {{- end}}

Use "{{.cmd.Name}} help --subcommand=[subcommand]" for more information about a subcommand.
`
	msg := cmd.executeTemplate(helpCommandTemplate, map[string]any{"cmd": cmd})

	return CommandResult{
		Title:      fmt.Sprintf("%v %v", cmd.Name, cmd.Emoji),
		Message:    msg,
		Successful: false,
	}
}

func (cmd *Command) HasAppID(appID entity.PlatformID) bool {
	return slices.Contains(cmd.AppIDs, appID)
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
		Name:       "help",
		Help:       fmt.Sprintf("Help for %v command", cmd.Name),
		AppIDs:     entity.AllAppIDs(),
		TargetFlag: TargetMaskAll,
		Handler: func(_ *entity.User, _ *Command, _ map[string]string) CommandResult {
			return cmd.RenderHelpTemplate()
		},
	}

	cmd.AddSubCommand(helpCmd)
}

func (cmd *Command) AddAboutSubCommand() {
	const aboutTemplate = `
## About Pagu

Version : {{.version}}
`

	cmd.ResultTemplate = aboutTemplate
	aboutCmd := &Command{
		Name:           "about",
		Help:           "About Pagu",
		AppIDs:         entity.AllAppIDs(),
		TargetFlag:     TargetMaskAll,
		ResultTemplate: aboutTemplate,
		Handler: func(_ *entity.User, _ *Command, _ map[string]string) CommandResult {
			return cmd.RenderResultTemplate("version", version.StringVersion())
		},
	}

	cmd.AddSubCommand(aboutCmd)
}
