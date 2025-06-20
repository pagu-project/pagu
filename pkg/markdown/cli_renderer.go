package markdown

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/pagu-project/pagu/pkg/log"
)

type CLIRenderer struct {
	renderer *glamour.TermRenderer
}

func NewCLIRenderer() *CLIRenderer {
	r, err := glamour.NewTermRenderer(
		glamour.WithColorProfile(lipgloss.ColorProfile()),
		glamour.WithAutoStyle(),
		glamour.WithPreservedNewLines(),
	)
	if err != nil {
		log.Warn("err in initialize Markdown renderer", "error", err)
	}

	return &CLIRenderer{
		renderer: r,
	}
}

func (md *CLIRenderer) Render(input string) string {
	if md.renderer == nil {
		return input
	}

	res, err := md.renderer.Render(input)
	if err != nil {
		log.Warn("error in rendering markdown", "error", err)

		return input
	}

	return res
}
