package variables

import (
	"strings"
	"unicode"
)

// ExpandFields replaces $VAR references in each field using values from store.
func ExpandFields(store *VariablesStore, fields []string) []string {
	if store == nil || len(fields) == 0 {
		return fields
	}

	expanded := make([]string, len(fields))
	for i, field := range fields {
		expanded[i] = ExpandField(store, field)
	}
	return expanded
}

// ExpandField replaces $VAR references in field with values from store.
// Undefined variables are left unchanged.
func ExpandField(store *VariablesStore, field string) string {
	if store == nil || !strings.Contains(field, "$") {
		return field
	}

	runes := []rune(field)
	var b strings.Builder

	for i := 0; i < len(runes); i++ {
		if runes[i] != '$' {
			b.WriteRune(runes[i])
			continue
		}
		if i+1 >= len(runes) {
			b.WriteRune('$')
			continue
		}

		name, length := readVariableName(runes[i+1:])
		if length == 0 {
			b.WriteRune('$')
			continue
		}

		value, ok := store.Get(name)
		if ok {
			b.WriteString(value)
		} else {
			b.WriteRune('$')
			b.WriteString(name)
		}
		i += length
	}

	return b.String()
}

func readVariableName(runes []rune) (string, int) {
	if len(runes) == 0 {
		return "", 0
	}
	if runes[0] != '_' && !unicode.IsLetter(runes[0]) {
		return "", 0
	}

	i := 1
	for i < len(runes) {
		r := runes[i]
		if r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) {
			i++
			continue
		}
		break
	}
	return string(runes[:i]), i
}
