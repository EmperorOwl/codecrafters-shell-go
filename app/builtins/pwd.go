package builtins

import (
	"fmt"
	"io"
	"os"
)

func init() {
	register("pwd", pwdHandler)
}

func pwdHandler(ctx *Context, args []string) (bool, error) {
	pwdBuiltin(ctx.Stdout)
	return false, nil
}

func pwdBuiltin(stdout io.Writer) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	fmt.Fprintln(stdout, cwd)
}
