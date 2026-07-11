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

// ExpandField replaces $VAR and ${VAR} references in field with values from store.
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
		i = expandAt(store, runes, i, &b)
	}

	return b.String()
}

func expandAt(store *VariablesStore, runes []rune, start int, b *strings.Builder) int {
	if start+1 >= len(runes) {
		b.WriteRune('$')
		return start
	}

	if runes[start+1] == '{' {
		closeIdx := findClosingBrace(runes, start+2)
		if closeIdx == -1 {
			b.WriteRune('$')
			return start
		}

		name := string(runes[start+2 : closeIdx])
		literal := "${" + name + "}"
		if IsValidIdentifier(name) {
			writeExpansion(b, store, name, literal)
		} else {
			b.WriteString(literal)
		}
		return closeIdx
	}

	name, length := readVariableName(runes[start+1:])
	if length == 0 {
		b.WriteRune('$')
		return start
	}

	writeExpansion(b, store, name, "$"+name)
	return start + length
}

func writeExpansion(b *strings.Builder, store *VariablesStore, name, literal string) {
	value, ok := store.Get(name)
	if ok {
		b.WriteString(value)
		return
	}
	b.WriteString(literal)
}

func findClosingBrace(runes []rune, start int) int {
	for i := start; i < len(runes); i++ {
		if runes[i] == '}' {
			return i
		}
	}
	return -1
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
