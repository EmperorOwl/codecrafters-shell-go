package terminal

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTerminal_WriteLine(t *testing.T) {
	tests := []struct {
		name    string
		rawMode bool
		text    string
		want    string
	}{
		{
			name:    "cooked mode appends newline",
			rawMode: false,
			text:    "hello",
			want:    "hello\n",
		},
		{
			name:    "raw mode translates newline",
			rawMode: true,
			text:    "hello",
			want:    "hello\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			term := New(nil, nil, strings.NewReader(""), &out, io.Discard)
			term.rawMode = tt.rawMode

			term.WriteLine(tt.text)

			if diff := cmp.Diff(tt.want, out.String()); diff != "" {
				t.Errorf("WriteLine(%q) mismatch (-want +got):\n%s", tt.text, diff)
			}
		})
	}
}

func TestTerminal_StdoutStderrDynamicWrapping(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*Terminal)
		write      func(*Terminal)
		wantStdout string
		wantStderr string
	}{
		{
			name: "cooked mode passes bytes through",
			setup: func(term *Terminal) {
				term.rawMode = false
			},
			write: func(term *Terminal) {
				term.Stdout().Write([]byte("out\n"))
				term.Stderr().Write([]byte("err\n"))
			},
			wantStdout: "out\n",
			wantStderr: "err\n",
		},
		{
			name: "raw mode translates on each write",
			setup: func(term *Terminal) {
				term.rawMode = true
			},
			write: func(term *Terminal) {
				term.Stdout().Write([]byte("out\n"))
				term.Stderr().Write([]byte("err\n"))
			},
			wantStdout: "out\r\n",
			wantStderr: "err\r\n",
		},
		{
			name: "wrapping follows raw mode changes",
			setup: func(term *Terminal) {
				term.rawMode = false
			},
			write: func(term *Terminal) {
				term.Stdout().Write([]byte("a\n"))
				term.rawMode = true
				term.Stdout().Write([]byte("b\n"))
				term.Stderr().Write([]byte("c\n"))
			},
			wantStdout: "a\nb\r\n",
			wantStderr: "c\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			term := New(nil, nil, strings.NewReader(""), &stdout, &stderr)
			if tt.setup != nil {
				tt.setup(term)
			}
			tt.write(term)

			if diff := cmp.Diff(tt.wantStdout, stdout.String()); diff != "" {
				t.Errorf("stdout mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantStderr, stderr.String()); diff != "" {
				t.Errorf("stderr mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestModeAwareWriter(t *testing.T) {
	tests := []struct {
		name       string
		stderr     bool
		setup      func(*Terminal)
		input      string
		wantStdout string
		wantStderr string
	}{
		{
			name:       "stdout in cooked mode",
			stderr:     false,
			setup:      func(term *Terminal) { term.rawMode = false },
			input:      "line\n",
			wantStdout: "line\n",
		},
		{
			name:       "stderr in raw mode",
			stderr:     true,
			setup:      func(term *Terminal) { term.rawMode = true },
			input:      "line\n",
			wantStderr: "line\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			term := New(nil, nil, strings.NewReader(""), &stdout, &stderr)
			if tt.setup != nil {
				tt.setup(term)
			}

			writer := modeAwareWriter{term: term, stderr: tt.stderr}
			if _, err := writer.Write([]byte(tt.input)); err != nil {
				t.Fatalf("Write() error = %v", err)
			}

			if diff := cmp.Diff(tt.wantStdout, stdout.String()); diff != "" {
				t.Errorf("stdout mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantStderr, stderr.String()); diff != "" {
				t.Errorf("stderr mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
