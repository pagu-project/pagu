package markdown

import (
	"log"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
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
		log.Printf("err in initialize Markdown renderer: %s", err)
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
		log.Printf("error in rendering markdown: %s", err)

		return input
	}

	return res
}
