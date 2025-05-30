package markdown_test

import (
	"testing"

	"github.com/pagu-project/pagu/pkg/markdown"
	"github.com/stretchr/testify/assert"
)

func TestMarkdownToHTML_Render(t *testing.T) {
	renderer := markdown.NewMarkdownToHTML()

	input := `# Hello

This is a [link](https://example.com) and some **bold** text.
`
	expected := `<h1>Hello</h1>

<p>This is a <a href="https://example.com" target="_blank">link</a> and some <strong>bold</strong> text.</p>
`

	output := renderer.Render(input)
	assert.Equal(t, expected, output)
}
