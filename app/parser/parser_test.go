package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  Command
	}{
		{
			name:  "simple command",
			input: []string{"echo", "hello"},
			want:  Command{Fields: []string{"echo", "hello"}},
		},
		{
			name:  "background",
			input: []string{"sleep", "1", "&"},
			want:  Command{Fields: []string{"sleep", "1"}, Background: true},
		},
		{
			name:  "stdout redirect",
			input: []string{"echo", "hi", ">", "out.txt"},
			want: Command{
				Fields:   []string{"echo", "hi"},
				Redirect: Redirect{StdoutPath: "out.txt"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseCommand(tt.input)
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ParseCommand() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParsePipelineSegments(t *testing.T) {
	tests := []struct {
		name      string
		segments  [][]string
		wantCmds  [][]string
		wantRedir Redirect
	}{
		{
			name: "strips background from segments",
			segments: [][]string{
				{"echo", "a"},
				{"cat", "&"},
			},
			wantCmds: [][]string{{"echo", "a"}, {"cat"}},
		},
		{
			name: "redirect on final segment",
			segments: [][]string{
				{"echo", "hello"},
				{"wc", ">", "out.txt"},
			},
			wantCmds:  [][]string{{"echo", "hello"}, {"wc"}},
			wantRedir: Redirect{StdoutPath: "out.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmds, gotRedir := ParsePipelineSegments(tt.segments)
			if diff := cmp.Diff(tt.wantCmds, gotCmds, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("commands mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantRedir, gotRedir); diff != "" {
				t.Errorf("redirect mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name string
		line string
		want Line
	}{
		{
			name: "single command",
			line: "echo hello",
			want: Line{Commands: [][]string{{"echo", "hello"}}},
		},
		{
			name: "background single command",
			line: "sleep 1 &",
			want: Line{
				Commands:   [][]string{{"sleep", "1"}},
				Background: true,
			},
		},
		{
			name: "pipeline",
			line: "echo hello | type echo",
			want: Line{
				Pipeline: true,
				Commands: [][]string{{"echo", "hello"}, {"type", "echo"}},
			},
		},
		{
			name: "pipeline with redirect",
			line: "echo hi | wc > out.txt",
			want: Line{
				Pipeline: true,
				Commands: [][]string{{"echo", "hi"}, {"wc"}},
				Redirect: Redirect{StdoutPath: "out.txt"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseLine(tt.line)
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ParseLine(%q) mismatch (-want +got):\n%s", tt.line, diff)
			}
		})
	}
}
