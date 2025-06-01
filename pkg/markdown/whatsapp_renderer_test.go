package markdown_test

import (
	"testing"

	"github.com/pagu-project/pagu/pkg/markdown"
	"github.com/stretchr/testify/assert"
)

func TestWhatsAppRenderer(t *testing.T) {
	renderer := markdown.NewWhatsAppRenderer()

	input := `# Hello
This is a [link](https://example.com) and some **bold** text.
`
	expected := `# Hello
This is a [link](https://example.com) and some **bold** text.
`

	output := renderer.Render(input)
	assert.Equal(t, expected, output)
}
