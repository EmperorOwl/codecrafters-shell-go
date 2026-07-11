package builtins

import (
	"fmt"
	"io"
)

func init() {
	register("declare", declareBuiltin)
}

func declareBuiltin(ctx *Context, args []string) (bool, error) {
	Declare(ctx.Stdout, ctx.Stderr, args)
	return false, nil
}

// Declare handles the declare builtin.
func Declare(stdout, stderr io.Writer, args []string) {
	if len(args) < 2 || args[0] != "-p" {
		return
	}
	fmt.Fprintln(stderr, variableNotFoundMessage(args[1]))
}

func variableNotFoundMessage(name string) string {
	return "declare: " + name + ": not found"
}
