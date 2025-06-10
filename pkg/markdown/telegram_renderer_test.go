package markdown_test

import (
	"testing"

	"github.com/pagu-project/pagu/pkg/markdown"
	"github.com/stretchr/testify/assert"
)

func TestTelegramUnescapeMarkdownLinks(t *testing.T) {
	renderer := markdown.NewTelegramRenderer()

	text := `Check out \[example\]\(https://example.com\) for more information.`
	expected := `Check out [example](https://example.com) for more information.`

	text = renderer.UnescapeMarkdownLinks(text)
	assert.Equal(t, expected, text)
}

func TestTelegramRender(t *testing.T) {
	renderer := markdown.NewTelegramRenderer()

	input := `# Hello
This is a link https://example.com and some **bold** and __italic__ text.
This is a [another link](https://example.com) with some (parentheses)[brackets].
## Header 2
”inline fixed-width code”
”  inline fixed-width code  ”
`
	expected := `*Hello*
This is a link https://example\.com and some *bold* and _italic_ text\.
This is a [another link](https://example\.com) with some \(parentheses\)\[brackets\]\.
*Header 2*
”inline fixed\-width code”
”  inline fixed\-width code  ”
`

	output := renderer.Render(input)
	assert.Equal(t, expected, output)
}
