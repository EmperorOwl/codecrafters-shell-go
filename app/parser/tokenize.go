package parser

import (
	"strings"
	"unicode"
)

const singleQuote = '\''

func Tokenize(line string) []string {
	var tokens []string
	var current strings.Builder
	inSingleQuote := false

	for _, r := range line {
		switch {
		case inSingleQuote:
			if r == singleQuote {
				inSingleQuote = false
			} else {
				current.WriteRune(r)
			}
		case r == singleQuote:
			inSingleQuote = true
		case unicode.IsSpace(r):
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}
