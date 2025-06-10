package markdown

import (
	"regexp"
	"strings"
)

type WhatsAppRenderer struct{}

func NewWhatsAppRenderer() *WhatsAppRenderer {
	return &WhatsAppRenderer{}
}

func (*WhatsAppRenderer) replaceMarkdownLinks(input string) string {
	// Regular expression to match [text](url)
	re := regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)

	return re.ReplaceAllString(input, "$1")
}

// Render converts the input string into a format compatible with WhatsApp.
// See: https://faq.whatsapp.com/539178204879377/?cms_platform=web
func (r *WhatsAppRenderer) Render(input string) string {
	// Headers: convert "# Heading" into "*Heading*"
	input = regexp.MustCompile(`(?m)^#+\s*(.+)$`).ReplaceAllString(input, "*$1*")

	// Bold: **text** → *text*
	input = strings.ReplaceAll(input, "**", "*")

	// Italic: __text__ → _text_
	input = strings.ReplaceAll(input, "__", "_")

	input = r.replaceMarkdownLinks(input)

	// Trim spaces inside quotes
	for {
		old := input
		input = strings.ReplaceAll(input, "” ", "”")
		input = strings.ReplaceAll(input, " ”", "”")
		if old == input {
			break
		}
	}

	return input
}
