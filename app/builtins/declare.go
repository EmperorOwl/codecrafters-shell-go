package builtins

import (
	"fmt"
	"io"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/variables"
)

func init() {
	register("declare", declareBuiltin)
}

func declareBuiltin(ctx *Context, args []string) (bool, error) {
	if ctx.State == nil {
		return false, nil
	}
	Declare(ctx.Stdout, ctx.Stderr, args, ctx.State.Variables)
	return false, nil
}

// Declare handles the declare builtin.
func Declare(stdout, stderr io.Writer, args []string, store *variables.VariablesStore) {
	if len(args) == 0 || store == nil {
		return
	}

	if args[0] == "-p" {
		if len(args) < 2 {
			return
		}
		name := args[1]
		value, ok := store.Get(name)
		if !ok {
			fmt.Fprintln(stderr, variableNotFoundMessage(name))
			return
		}
		fmt.Fprintln(stdout, variableDescriptionMessage(name, value))
		return
	}

	if name, value, ok := parseAssignment(args[0]); ok {
		store.Set(name, value)
	}
}

func parseAssignment(arg string) (name, value string, ok bool) {
	index := strings.Index(arg, "=")
	if index <= 0 {
		return "", "", false
	}
	return arg[:index], arg[index+1:], true
}

func variableDescriptionMessage(name, value string) string {
	return fmt.Sprintf(`declare -- %s="%s"`, name, value)
}

func variableNotFoundMessage(name string) string {
	return "declare: " + name + ": not found"
}
