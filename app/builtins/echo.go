package builtins

import (
	"fmt"
	"io"
	"strings"
)

func init() {
	register("echo", echoHandler)
}

func echoHandler(ctx *Context, args []string) (bool, error) {
	echoBuiltin(ctx.Stdout, args)
	return false, nil
}

func echoBuiltin(stdout io.Writer, args []string) {
	fmt.Fprintln(stdout, strings.Join(args, " "))
}
