package builtins

import (
	"fmt"
	"io"
	"strings"
)

func init() {
	register("echo", echoBuiltin)
}

func echoBuiltin(ctx *Context, args []string) (bool, error) {
	Echo(ctx.Stdout, args)
	return false, nil
}

func Echo(out io.Writer, args []string) {
	fmt.Fprintln(out, strings.Join(args, " "))
}
