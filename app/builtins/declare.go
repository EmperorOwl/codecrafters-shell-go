package builtins

func init() {
	register("declare", declareBuiltin)
}

func declareBuiltin(ctx *Context, args []string) (bool, error) {
	return false, nil
}
