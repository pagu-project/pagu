package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func CamelCase(input string) string {
	input = strings.ReplaceAll(input, "-", " ")
	words := strings.Fields(input)

	for i, word := range words {
		// Lowercase the first word, capitalize the rest
		if i == 0 {
			words[i] = strings.ToLower(word)
		} else {
			words[i] = cases.Title(language.English).String(word)
		}
	}

	return strings.Join(words, "")
}
