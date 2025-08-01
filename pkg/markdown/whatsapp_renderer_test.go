package markdown_test

import (
	"testing"

	"github.com/pagu-project/pagu/pkg/markdown"
	"github.com/stretchr/testify/assert"
)

func TestWhatsAppRenderer(t *testing.T) {
	renderer := markdown.NewWhatsAppRenderer()

	input := `# Hello
This is a link https://example.com and some **bold** and __italic__ text.
This is a [another link](https://example.com) with some (parentheses)[brackets].
## Header 2
”inline fixed-width code”
”  inline fixed-width code  ”
`
	expected := `*Hello*
This is a link https://example.com and some *bold* and _italic_ text.
This is a https://example.com with some (parentheses)[brackets].
*Header 2*
”inline fixed-width code”
”inline fixed-width code”
`

	output := renderer.Render(input)
	assert.Equal(t, expected, output)
}
