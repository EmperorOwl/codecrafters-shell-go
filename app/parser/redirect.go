package parser

const (
	stdoutRedirectOp           = ">"
	stdoutRedirectOpLong       = "1>"
	stdoutAppendRedirectOp     = ">>"
	stdoutAppendRedirectOpLong = "1>>"
	stderrRedirectOp           = "2>"
)

// ParseRedirect splits tokenized command arguments from optional stdout and stderr redirects.
// It recognizes >, 1>, >>, 1>>, and 2> followed by a path token.
func ParseRedirect(tokens []string) (fields []string, stdoutPath string, stdoutAppend bool, stderrPath string) {
	for i := 0; i < len(tokens); i++ {
		switch tokens[i] {
		case stdoutAppendRedirectOp, stdoutAppendRedirectOpLong:
			if i+1 < len(tokens) {
				stdoutPath = tokens[i+1]
				stdoutAppend = true
				i++
			}
		case stdoutRedirectOp, stdoutRedirectOpLong:
			if i+1 < len(tokens) {
				stdoutPath = tokens[i+1]
				i++
			}
		case stderrRedirectOp:
			if i+1 < len(tokens) {
				stderrPath = tokens[i+1]
				i++
			}
		default:
			fields = append(fields, tokens[i])
		}
	}
	return fields, stdoutPath, stdoutAppend, stderrPath
}
