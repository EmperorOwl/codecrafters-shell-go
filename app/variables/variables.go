package variables

// VariablesStore holds shell variable names and values.
type VariablesStore struct {
	values map[string]string
}

// NewVariablesStore returns an empty variable store.
func NewVariablesStore() *VariablesStore {
	return &VariablesStore{
		values: make(map[string]string),
	}
}

// Set records a shell variable value.
func (s *VariablesStore) Set(name, value string) {
	s.values[name] = value
}

// Get returns the value for a shell variable.
func (s *VariablesStore) Get(name string) (string, bool) {
	value, ok := s.values[name]
	return value, ok
}
