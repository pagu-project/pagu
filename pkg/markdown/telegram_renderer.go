package markdown

import (
	"regexp"
	"strings"
)

type TelegramRenderer struct{}

func NewTelegramRenderer() *TelegramRenderer {
	return &TelegramRenderer{}
}

func (*TelegramRenderer) escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
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
	)

	return replacer.Replace(text)
}

func (*TelegramRenderer) encodeMarkdownLinks(input string) string {
	re := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

	return re.ReplaceAllString(input, "†$1†‡$2‡")
}

func (*TelegramRenderer) decodeMarkdownLinks(input string) string {
	re := regexp.MustCompile(`†([^†]+)†‡([^‡]+)‡`)

	return re.ReplaceAllString(input, "[$1]($2)")
}

// Render converts the input string into a format compatible with Telegram's MarkdownV2 style.
// For formatting rules, see: https://core.telegram.org/bots/api#markdownv2-style
func (r *TelegramRenderer) Render(input string) string {
	// Headers: convert "# Heading" into "*Heading*"
	input = regexp.MustCompile(`(?m)^#+\s*(.+)$`).ReplaceAllString(input, "*$1*")

	// Bold: **text** → *text*
	input = strings.ReplaceAll(input, "**", "*")

	// Italic: __text__ → _text_
	input = strings.ReplaceAll(input, "__", "_")

	// [Google](https://google.com) → †Google†‡https://google.com‡
	input = r.encodeMarkdownLinks(input)

	// Escape remaining MarkdownV2 special characters
	input = r.escapeMarkdownV2(input)

	// †Google†‡https://google.com‡ → [Google](https://google.com)
	input = r.decodeMarkdownLinks(input)

	return input
}
