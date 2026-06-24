package shell

import (
	"fmt"
	"io"
	"strings"
)

func EchoOutput(args []string) string {
	return strings.Join(args, " ")
}

func TryBuiltin(line string, out io.Writer) (handled bool, shouldExit bool) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return false, false
	}

	switch fields[0] {
	case "exit":
		return true, true
	case "echo":
		fmt.Fprintln(out, EchoOutput(fields[1:]))
		return true, false
	default:
		return false, false
	}
}
