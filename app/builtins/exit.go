package builtins

func init() {
	register("exit", exitBuiltin)
}

func exitBuiltin(ctx *Context, args []string) (bool, error) {
	return true, nil
}
