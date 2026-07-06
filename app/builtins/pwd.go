package builtins

import (
	"fmt"
	"io"
	"os"
)

func PwdOutput() (string, error) {
	return os.Getwd()
}

func Pwd(out io.Writer) {
	cwd, err := PwdOutput()
	if err != nil {
		return
	}
	fmt.Fprintln(out, cwd)
}

func init() {
	register("pwd", pwdBuiltin)
}

func pwdBuiltin(ctx *Context, args []string) (bool, error) {
	Pwd(ctx.Stdout)
	return false, nil
}
