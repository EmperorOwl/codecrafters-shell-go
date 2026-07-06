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

func init() {
	register("echo", echoBuiltin)
}

func echoBuiltin(ctx *Context, args []string) (bool, error) {
	Echo(ctx.Stdout, args)
	return false, nil
}
