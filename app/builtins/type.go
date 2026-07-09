package builtins

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/external"
)

func init() {
	register("type", typeBuiltin)
}

func typeBuiltin(ctx *Context, args []string) (bool, error) {
	target := ""
	if len(args) > 0 {
		target = args[0]
	}
	Type(ctx.Stdout, target)
	return false, nil
}

func Type(out io.Writer, command string) {
	if IsBuiltin(command) {
		fmt.Fprintln(out, command+" is a shell builtin")
		return
	}
	if path, ok := external.FindExecutableInPath(command); ok {
		fmt.Fprintln(out, command+" is "+path)
		return
	}
	fmt.Fprintln(out, command+": not found")
}
