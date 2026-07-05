package completion

import (
	"slices"
	"strings"
)

// Complete finds matches for prefix among candidates. The bool is true when
// exactly one candidate matched.
func Complete(prefix string, candidates []string) (token string, listings []string, unique bool) {
	matched := findMatches(candidates, prefix)

	switch len(matched) {
	case 0:
		return prefix, nil, false
	case 1:
		return matched[0], nil, true
	default:
		lcp := longestCommonPrefix(matched)
		if len(lcp) > len(prefix) {
			return lcp, nil, false
		}
		return prefix, matched, false
	}
}

func findMatches(candidates []string, prefix string) []string {
	var result []string
	for _, candidate := range candidates {
		if strings.HasPrefix(candidate, prefix) {
			result = append(result, candidate)
		}
	}
	slices.Sort(result)
	return slices.Compact(result)
}

func longestCommonPrefix(candidates []string) string {
	if len(candidates) == 0 {
		return ""
	}

	prefix := candidates[0]
	for _, candidate := range candidates[1:] {
		for len(prefix) > 0 && !strings.HasPrefix(candidate, prefix) {
			prefix = prefix[:len(prefix)-1]
		}
		if prefix == "" {
			return ""
		}
	}

	return prefix
}
