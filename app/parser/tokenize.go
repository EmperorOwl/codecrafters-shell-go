package parser

import (
	"slices"
	"strings"
	"unicode"
)

const (
	singleQuote = '\''
	doubleQuote = '"'
	backslash   = '\\'
)

var backslashInDoubleQuotesEscapes = []rune{
	doubleQuote, backslash, '$', '`', '\n',
}

func Tokenize(line string) []string {
	var tokens []string
	var current strings.Builder
	inDoubleQuotes := false
	inSingleQuotes := false
	escaping := false

	handleTokenDone := func() {
		if current.Len() == 0 {
			return
		}
		tokens = append(tokens, current.String())
		current.Reset()
		inDoubleQuotes = false
		inSingleQuotes = false
		escaping = false
	}

	for _, r := range line {
		if escaping {
			if inDoubleQuotes && !slices.Contains(backslashInDoubleQuotesEscapes, r) {
				current.WriteRune(backslash)
			}
			current.WriteRune(r)
			escaping = false
			continue
		}

		switch {
		case r == doubleQuote && !inSingleQuotes:
			inDoubleQuotes = !inDoubleQuotes

		case r == singleQuote && !inDoubleQuotes:
			inSingleQuotes = !inSingleQuotes

		case r == backslash && !inSingleQuotes:
			escaping = true

		case unicode.IsSpace(r) && !inDoubleQuotes && !inSingleQuotes:
			handleTokenDone()

		default:
			current.WriteRune(r)
		}
	}

	handleTokenDone()
	return tokens
}
