package markdown

import (
	"log"

	"github.com/charmbracelet/glamour"
)

type Markdown struct {
	renderer *glamour.TermRenderer
}

type MarkdownInterface interface {
	Render(in string) (string, error)
}

func NewMarkdown() MarkdownInterface {
	r, err := glamour.NewTermRenderer(
		// glamour.WithColorProfile(lipgloss.ColorProfile()),
		// glamour.WithAutoStyle(),
		glamour.WithPreservedNewLines(),
	)

	if err != nil {
		log.Printf("err in ininitial mark down renderer: %s", err)
	}

	return &Markdown{
		renderer: r,
	}
}

func (md *Markdown) Render(in string) (string, error) {
	return md.renderer.Render(in)
}
