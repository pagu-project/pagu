package utils

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func CamelCase(input string) string {
	if input == "" {
		return ""
	}

	re := regexp.MustCompile(`[^a-zA-Z0-9\-_]+`)
	input = re.ReplaceAllString(input, " ")

	input = strings.NewReplacer("-", " ", "_", " ").Replace(input)

	words := strings.Fields(input)

	if len(words) == 0 {
		return ""
	}

	builder := strings.Builder{}
	caser := cases.Title(language.English)

	for i, word := range words {
		if i == 0 {
			builder.WriteString(strings.ToLower(word))
		} else {
			builder.WriteString(caser.String(word))
		}
	}

	return builder.String()
}
