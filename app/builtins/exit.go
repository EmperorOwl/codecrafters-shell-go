package builtins

func init() {
	register("exit", exitHandler)
}

func exitHandler(ctx *Context, args []string) (bool, error) {
	return true, nil
}
