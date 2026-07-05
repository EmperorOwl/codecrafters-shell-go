package completion

// CompletionRegistry stores programmable completion script paths per command.
type CompletionRegistry struct {
	scripts map[string]string
}

// NewCompletionRegistry returns an empty completion registry.
func NewCompletionRegistry() *CompletionRegistry {
	return &CompletionRegistry{
		scripts: make(map[string]string),
	}
}

// Register associates a completer script with a command.
func (r *CompletionRegistry) Register(command, scriptPath string) {
	r.scripts[command] = scriptPath
}

// Unregister removes the completion rule for a command.
func (r *CompletionRegistry) Unregister(command string) {
	delete(r.scripts, command)
}

// Lookup returns the completer script path for a command.
func (r *CompletionRegistry) Lookup(command string) (string, bool) {
	scriptPath, ok := r.scripts[command]
	return scriptPath, ok
}
