package completion

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// CompleterOptions holds the context passed to a programmable completer script.
type CompleterOptions struct {
	Path         string
	Command      string
	CurrentWord  string
	PreviousWord string
	CompLine     string
	CompPoint    int
}

// RunCompleter executes a programmable completer script and returns its stdout
// as a list of candidate strings. Path on opts must be set.
func RunCompleter(opts CompleterOptions) ([]string, error) {
	cmd := exec.Command(opts.Path, opts.Command, opts.CurrentWord, opts.PreviousWord)
	cmd.Env = append(os.Environ(),
		"COMP_LINE="+opts.CompLine,
		"COMP_POINT="+strconv.Itoa(opts.CompPoint),
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseCompleterOutput(output), nil
}

// parseCompleterOutput splits completer script stdout into one candidate per line.
func parseCompleterOutput(output []byte) []string {
	text := strings.TrimRight(string(output), "\r\n")
	if text == "" {
		return nil
	}

	lines := strings.Split(text, "\n")
	candidates := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSuffix(line, "\r")
		if line != "" {
			candidates = append(candidates, line)
		}
	}
	return candidates
}
