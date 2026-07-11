package variables

import "unicode"

// IsValidIdentifier reports whether name is a valid shell variable name.
func IsValidIdentifier(name string) bool {
	if name == "" {
		return false
	}
	for i, r := range name {
		switch {
		case i == 0:
			if r != '_' && !unicode.IsLetter(r) {
				return false
			}
		case r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r):
			return false
		}
	}
	return true
}
