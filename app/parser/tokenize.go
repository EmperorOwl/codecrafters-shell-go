package parser

import (
	"strings"
	"unicode"
)

const (
	singleQuote = '\''
	doubleQuote = '"'
	backslash   = '\\'
)

func Tokenize(line string) []string {
	runes := []rune(line)
	var tokens []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case inSingleQuote:
			if r == singleQuote {
				inSingleQuote = false
			} else {
				current.WriteRune(r)
			}
		case inDoubleQuote:
			if r == backslash {
				if i+1 < len(runes) {
					next := runes[i+1]
					switch next {
					case doubleQuote, backslash:
						i++
						current.WriteRune(next)
					default:
						current.WriteRune(r)
					}
				} else {
					current.WriteRune(r)
				}
			} else if r == doubleQuote {
				inDoubleQuote = false
			} else {
				current.WriteRune(r)
			}
		case r == backslash:
			if i+1 < len(runes) {
				i++
				current.WriteRune(runes[i])
			}
		case r == singleQuote:
			inSingleQuote = true
		case r == doubleQuote:
			inDoubleQuote = true
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
