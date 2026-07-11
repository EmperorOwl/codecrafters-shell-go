package completion

// Registry stores programmable completion script paths per command.
type Registry struct {
	scripts map[string]string
}

// NewRegistry returns an empty completion registry.
func NewRegistry() *Registry {
	return &Registry{
		scripts: make(map[string]string),
	}
}

// Register associates a completer script with a command.
func (r *Registry) Register(command, scriptPath string) {
	r.scripts[command] = scriptPath
}

// Unregister removes the completion rule for a command.
func (r *Registry) Unregister(command string) {
	delete(r.scripts, command)
}

// Lookup returns the completer script path for a command.
func (r *Registry) Lookup(command string) (string, bool) {
	scriptPath, ok := r.scripts[command]
	return scriptPath, ok
}
