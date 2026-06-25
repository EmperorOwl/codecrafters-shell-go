package parser

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "empty", input: "", want: nil},
		{name: "single word", input: "echo", want: []string{"echo"}},
		{
			name:  "spaces preserved in single quotes",
			input: "echo 'hello    world'",
			want:  []string{"echo", "hello    world"},
		},
		{
			name:  "unquoted spaces collapse between tokens",
			input: "echo hello    world",
			want:  []string{"echo", "hello", "world"},
		},
		{
			name:  "adjacent quoted strings concatenate",
			input: "echo 'hello''world'",
			want:  []string{"echo", "helloworld"},
		},
		{
			name:  "empty quotes are ignored",
			input: "echo hello''world",
			want:  []string{"echo", "helloworld"},
		},
		{
			name:  "quoted file paths for cat",
			input: "cat '/tmp/file name' '/tmp/file name with spaces'",
			want:  []string{"cat", "/tmp/file name", "/tmp/file name with spaces"},
		},
		{
			name:  "special characters literal in quotes",
			input: "echo '$HOME * ~'",
			want:  []string{"echo", "$HOME * ~"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Tokenize(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokenize(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
