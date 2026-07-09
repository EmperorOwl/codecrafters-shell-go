package builtins

import (
	"fmt"
	"io"
	"os"
)

func init() {
	register("pwd", pwdBuiltin)
}

func pwdBuiltin(ctx *Context, args []string) (bool, error) {
	Pwd(ctx.Stdout)
	return false, nil
}

func Pwd(out io.Writer) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	fmt.Fprintln(out, cwd)
}
