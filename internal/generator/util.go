package generator

import (
	"strings"
	"unicode"
)

func snakeToCamel(input string) string {
	// Handle empty string
	if len(input) == 0 {
		return input
	}

	// Split the input string by underscores
	words := strings.Split(input, "_")

	// Process all words, including the first one
	var result strings.Builder
	for _, word := range words {
		if len(word) > 0 {
			// Capitalize the first letter of each word
			firstChar := unicode.ToUpper(rune(word[0]))
			result.WriteRune(firstChar)

			if len(word) > 1 {
				// Lowercase the rest of the word
				result.WriteString(strings.ToLower(word[1:]))
			}
		}
	}

	return result.String()
}
