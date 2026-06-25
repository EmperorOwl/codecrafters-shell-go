package parser

const (
	stdoutRedirectOp           = ">"
	stdoutRedirectOpLong       = "1>"
	stdoutAppendRedirectOp     = ">>"
	stdoutAppendRedirectOpLong = "1>>"
	stderrRedirectOp           = "2>"
	stderrAppendRedirectOp     = "2>>"
)

type Redirect struct {
	StdoutPath   string
	StdoutAppend bool
	StderrPath   string
	StderrAppend bool
}

// ParseRedirect splits tokenized command arguments from optional stdout and stderr redirects.
// It recognizes >, 1>, >>, 1>>, 2>, and 2>> followed by a path token.
func ParseRedirect(tokens []string) (fields []string, redirect Redirect) {
	for i := 0; i < len(tokens); i++ {
		switch tokens[i] {
		case stdoutAppendRedirectOp, stdoutAppendRedirectOpLong:
			if i+1 < len(tokens) {
				redirect.StdoutPath = tokens[i+1]
				redirect.StdoutAppend = true
				i++
			}
		case stdoutRedirectOp, stdoutRedirectOpLong:
			if i+1 < len(tokens) {
				redirect.StdoutPath = tokens[i+1]
				i++
			}
		case stderrAppendRedirectOp:
			if i+1 < len(tokens) {
				redirect.StderrPath = tokens[i+1]
				redirect.StderrAppend = true
				i++
			}
		case stderrRedirectOp:
			if i+1 < len(tokens) {
				redirect.StderrPath = tokens[i+1]
				i++
			}
		default:
			fields = append(fields, tokens[i])
		}
	}
	return fields, redirect
}
