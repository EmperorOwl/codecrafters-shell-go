package shell

func HandleBuiltin(command string) bool {
	switch command {
	case "exit":
		return true
	default:
		return false
	}
}
