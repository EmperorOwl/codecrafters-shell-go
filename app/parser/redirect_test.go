package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseRedirect(t *testing.T) {
	tests := []struct {
		name         string
		tokens       []string
		wantFields   []string
		wantRedirect Redirect
	}{
		{
			name:       "no redirect",
			tokens:     []string{"echo", "hello"},
			wantFields: []string{"echo", "hello"},
		},
		{
			name:         "stdout redirect",
			tokens:       []string{"echo", "hello", ">", "output.txt"},
			wantFields:   []string{"echo", "hello"},
			wantRedirect: Redirect{StdoutPath: "output.txt"},
		},
		{
			name:         "explicit stdout redirect",
			tokens:       []string{"echo", "Hello", "James", "1>", "/tmp/foo/foo.md"},
			wantFields:   []string{"echo", "Hello", "James"},
			wantRedirect: Redirect{StdoutPath: "/tmp/foo/foo.md"},
		},
		{
			name:         "external command with redirect",
			tokens:       []string{"ls", "/tmp/baz", ">", "/tmp/foo/baz.md"},
			wantFields:   []string{"ls", "/tmp/baz"},
			wantRedirect: Redirect{StdoutPath: "/tmp/foo/baz.md"},
		},
		{
			name:         "redirect with multiple arguments",
			tokens:       []string{"cat", "/tmp/baz/blueberry", "nonexistent", "1>", "/tmp/foo/quz.md"},
			wantFields:   []string{"cat", "/tmp/baz/blueberry", "nonexistent"},
			wantRedirect: Redirect{StdoutPath: "/tmp/foo/quz.md"},
		},
		{
			name:       "greater than inside token is not redirect",
			tokens:     []string{"echo", "a>b"},
			wantFields: []string{"echo", "a>b"},
		},
		{
			name:         "stderr redirect",
			tokens:       []string{"ls", "nonexistent", "2>", "/tmp/quz/baz.md"},
			wantFields:   []string{"ls", "nonexistent"},
			wantRedirect: Redirect{StderrPath: "/tmp/quz/baz.md"},
		},
		{
			name:         "stderr redirect with stdout output",
			tokens:       []string{"cat", "/tmp/bar/pear", "nonexistent", "2>", "/tmp/quz/quz.md"},
			wantFields:   []string{"cat", "/tmp/bar/pear", "nonexistent"},
			wantRedirect: Redirect{StderrPath: "/tmp/quz/quz.md"},
		},
		{
			name:       "stdout append redirect",
			tokens:     []string{"echo", "Hello", "Emily", "1>>", "/tmp/bar/baz.md"},
			wantFields: []string{"echo", "Hello", "Emily"},
			wantRedirect: Redirect{
				StdoutPath:   "/tmp/bar/baz.md",
				StdoutAppend: true,
			},
		},
		{
			name:       "external command with append redirect",
			tokens:     []string{"ls", "/tmp/baz", ">>", "/tmp/bar/bar.md"},
			wantFields: []string{"ls", "/tmp/baz"},
			wantRedirect: Redirect{
				StdoutPath:   "/tmp/bar/bar.md",
				StdoutAppend: true,
			},
		},
		{
			name:       "double greater than inside token is not redirect",
			tokens:     []string{"echo", "a>>b"},
			wantFields: []string{"echo", "a>>b"},
		},
		{
			name:       "stderr append redirect",
			tokens:     []string{"ls", "nonexistent", "2>>", "/tmp/foo/qux.md"},
			wantFields: []string{"ls", "nonexistent"},
			wantRedirect: Redirect{
				StderrPath:   "/tmp/foo/qux.md",
				StderrAppend: true,
			},
		},
		{
			name:       "double stderr redirect inside token is not redirect",
			tokens:     []string{"echo", "a2>>b"},
			wantFields: []string{"echo", "a2>>b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFields, gotRedirect := ParseRedirect(tt.tokens)
			if diff := cmp.Diff(tt.wantFields, gotFields); diff != "" {
				t.Errorf("ParseRedirect() fields mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantRedirect, gotRedirect); diff != "" {
				t.Errorf("ParseRedirect() redirect mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
