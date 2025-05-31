package markdown

import (
	"regexp"
	"strings"
)

type TelegramRenderer struct{}

func NewTelegramRenderer() *TelegramRenderer {
	return &TelegramRenderer{}
}

func escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		`\,`, `\\`,
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)

	return replacer.Replace(text)
}

// Render converts the input string into a format compatible with Telegram's MarkdownV2 style.
// For formatting rules, see: https://core.telegram.org/bots/api#markdownv2-style
func (*TelegramRenderer) Render(input string) string {
	// Headers: convert "# Heading" into "*Heading*"
	input = regexp.MustCompile(`(?m)^#+\s*(.+)$`).ReplaceAllString(input, "*$1*")

	// Bold: **text** → *text*
	input = strings.ReplaceAll(input, "**", "*")

	// Italic: __text__ → _text_
	input = strings.ReplaceAll(input, "__", "_")

	// Escape remaining MarkdownV2 special characters
	return escapeMarkdownV2(input)
}
