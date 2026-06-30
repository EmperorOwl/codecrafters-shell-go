package parser

// StripBackground removes a trailing & token and reports whether the command
// should run in the background.
func StripBackground(tokens []string) (fields []string, background bool) {
	if len(tokens) == 0 || tokens[len(tokens)-1] != "&" {
		return tokens, false
	}
	return tokens[:len(tokens)-1], true
}
