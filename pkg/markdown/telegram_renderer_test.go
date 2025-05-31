package markdown_test

import (
	"testing"

	"github.com/pagu-project/pagu/pkg/markdown"
	"github.com/stretchr/testify/assert"
)

func TestTelegramRender(t *testing.T) {
	renderer := markdown.NewTelegramRenderer()

	input := `# Hello
This is a link https://example.com and some **bold** and __italic__ text.
”inline fixed-width code”
`
	expected := `*Hello*
This is a link https://example\.com and some *bold* and _italic_ text\.
”inline fixed\-width code”
`

	output := renderer.Render(input)
	assert.Equal(t, expected, output)
}
