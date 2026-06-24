package builtins

import (
	"fmt"
	"io"
	"strings"
)

func EchoOutput(args []string) string {
	return strings.Join(args, " ")
}

func Echo(out io.Writer, args []string) {
	fmt.Fprintln(out, EchoOutput(args))
}
