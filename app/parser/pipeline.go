package parser

// SplitPipelineTokens splits tokenized arguments on | operators.
// Returns nil when the tokens contain no pipeline operator.
func SplitPipelineTokens(tokens []string) [][]string {
	var segments [][]string
	var current []string
	foundPipe := false

	for _, token := range tokens {
		if token == string(pipeOp) {
			foundPipe = true
			segments = append(segments, current)
			current = nil
			continue
		}
		current = append(current, token)
	}

	if !foundPipe {
		return nil
	}

	segments = append(segments, current)
	return segments
}
