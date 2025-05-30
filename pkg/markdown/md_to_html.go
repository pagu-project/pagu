package markdown

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type MarkdownToHTML struct {
	parser *parser.Parser
}

func NewMarkdownToHTML() *MarkdownToHTML {
	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
	parser := parser.NewWithExtensions(extensions)

	return &MarkdownToHTML{
		parser: parser,
	}
}

func (md *MarkdownToHTML) Render(input string) string {
	doc := md.parser.Parse([]byte(input))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	html := markdown.Render(doc, renderer)

	return string(html)
}
