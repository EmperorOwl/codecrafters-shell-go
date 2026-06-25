package parser

const (
	stdoutRedirectOp     = ">"
	stdoutRedirectOpLong = "1>"
)

// ParseRedirect splits tokenized command arguments from an optional stdout redirect.
// It recognizes > and 1> followed by a path token.
func ParseRedirect(tokens []string) (fields []string, stdoutPath string) {
	for i := 0; i < len(tokens); i++ {
		switch tokens[i] {
		case stdoutRedirectOp, stdoutRedirectOpLong:
			if i+1 < len(tokens) {
				stdoutPath = tokens[i+1]
				i++
			}
		default:
			fields = append(fields, tokens[i])
		}
	}
	return fields, stdoutPath
}
