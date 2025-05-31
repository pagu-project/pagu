package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
)

type MarkdownToHTML struct {
	renderer goldmark.Markdown
}

func NewMarkdownToHTML() *MarkdownToHTML {
	renderer := goldmark.New()

	return &MarkdownToHTML{
		renderer: renderer,
	}
}

func (md *MarkdownToHTML) Render(input string) string {
	var buf bytes.Buffer
	if err := md.renderer.Convert([]byte(input), &buf); err != nil {
		panic(err)
	}

	return buf.String()
}
