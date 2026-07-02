package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSplitPipelineTokens(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  [][]string
	}{
		{name: "no pipeline", input: "echo hello", want: nil},
		{
			name:  "two commands",
			input: "cat /tmp/foo/file | wc",
			want:  [][]string{{"cat", "/tmp/foo/file"}, {"wc"}},
		},
		{
			name:  "pipe with surrounding spaces",
			input: "tail -f /tmp/foo/file-1 | head -n 5",
			want:  [][]string{{"tail", "-f", "/tmp/foo/file-1"}, {"head", "-n", "5"}},
		},
		{
			name:  "pipe without surrounding spaces",
			input: "cat file|wc",
			want:  [][]string{{"cat", "file"}, {"wc"}},
		},
		{name: "pipe inside single quotes", input: "echo 'a|b'", want: nil},
		{name: "pipe inside double quotes", input: `echo "a|b"`, want: nil},
		{name: "escaped pipe is literal", input: `echo a\|b`, want: nil},
		{
			name:  "three commands",
			input: "a | b | c",
			want:  [][]string{{"a"}, {"b"}, {"c"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitPipelineTokens(Tokenize(tt.input))
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("SplitPipelineTokens(Tokenize(%q)) mismatch (-want +got):\n%s", tt.input, diff)
			}
		})
	}
}
