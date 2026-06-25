package parser

import (
	"reflect"
	"testing"
)

func TestParseRedirect(t *testing.T) {
	tests := []struct {
		name           string
		tokens         []string
		wantFields     []string
		wantStdoutPath string
	}{
		{
			name:       "no redirect",
			tokens:     []string{"echo", "hello"},
			wantFields: []string{"echo", "hello"},
		},
		{
			name:           "stdout redirect",
			tokens:         []string{"echo", "hello", ">", "output.txt"},
			wantFields:     []string{"echo", "hello"},
			wantStdoutPath: "output.txt",
		},
		{
			name:           "explicit stdout redirect",
			tokens:         []string{"echo", "Hello", "James", "1>", "/tmp/foo/foo.md"},
			wantFields:     []string{"echo", "Hello", "James"},
			wantStdoutPath: "/tmp/foo/foo.md",
		},
		{
			name:           "external command with redirect",
			tokens:         []string{"ls", "/tmp/baz", ">", "/tmp/foo/baz.md"},
			wantFields:     []string{"ls", "/tmp/baz"},
			wantStdoutPath: "/tmp/foo/baz.md",
		},
		{
			name:           "redirect with multiple arguments",
			tokens:         []string{"cat", "/tmp/baz/blueberry", "nonexistent", "1>", "/tmp/foo/quz.md"},
			wantFields:     []string{"cat", "/tmp/baz/blueberry", "nonexistent"},
			wantStdoutPath: "/tmp/foo/quz.md",
		},
		{
			name:       "greater than inside token is not redirect",
			tokens:     []string{"echo", "a>b"},
			wantFields: []string{"echo", "a>b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFields, gotStdoutPath := ParseRedirect(tt.tokens)
			if !reflect.DeepEqual(gotFields, tt.wantFields) {
				t.Errorf("ParseRedirect() fields = %v, want %v", gotFields, tt.wantFields)
			}
			if gotStdoutPath != tt.wantStdoutPath {
				t.Errorf("ParseRedirect() stdoutPath = %q, want %q", gotStdoutPath, tt.wantStdoutPath)
			}
		})
	}
}
