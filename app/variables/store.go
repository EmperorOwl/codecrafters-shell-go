package variables

// Store holds shell variable names and values.
type Store struct {
	values map[string]string
}

// NewStore returns an empty variable store.
func NewStore() *Store {
	return &Store{
		values: make(map[string]string),
	}
}

// Set records a shell variable value.
func (s *Store) Set(name, value string) {
	s.values[name] = value
}

// Get returns the value for a shell variable.
func (s *Store) Get(name string) (string, bool) {
	value, ok := s.values[name]
	return value, ok
}

// Entries returns a copy of all stored variables.
func (s *Store) Entries() map[string]string {
	entries := make(map[string]string, len(s.values))
	for name, value := range s.values {
		entries[name] = value
	}
	return entries
}
