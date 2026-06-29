package completion

import (
	"slices"
)

func applyCommandTab(builtins, executables []string, buffer string) (newBuffer string, listings []string) {
	matched := append(findMatches(builtins, buffer), findMatches(executables, buffer)...)
	slices.Sort(matched)
	matched = slices.Compact(matched)

	switch len(matched) {
	case 0:
		return buffer, nil
	case 1:
		return matched[0] + " ", nil
	default:
		lcp := longestCommonPrefix(matched)
		if len(lcp) > len(buffer) {
			return lcp, nil
		}
		return buffer, matched
	}
}
