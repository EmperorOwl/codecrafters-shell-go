package shell

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRunWritesHistfileOnExit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "histfile")
	t.Setenv("HISTFILE", path)

	input := strings.NewReader("echo hello\necho world\nexit\n")
	var out bytes.Buffer
	s := New(input, &out, io.Discard)

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	want := "echo hello\necho world\nexit\n"
	if diff := cmp.Diff(want, string(got)); diff != "" {
		t.Errorf("Run() histfile mismatch (-want +got):\n%s", diff)
	}
}

func TestCommandNotFoundMessage(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    string
	}{
		{
			name:    "unknown command",
			command: "missing",
			want:    "missing: command not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CommandNotFoundMessage(tt.command)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("CommandNotFoundMessage(%q) mismatch (-want +got):\n%s", tt.command, diff)
			}
		})
	}
}

func TestWriteReapedJobs(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*jobs.Table)
		wantLines []string
	}{
		{
			name:      "no done jobs",
			setup:     func(*jobs.Table) {},
			wantLines: nil,
		},
		{
			name: "running job produces no output",
			setup: func(jm *jobs.Table) {
				jm.Add(1, "sleep 10 &")
			},
			wantLines: nil,
		},
		{
			name: "prints one done job",
			setup: func(jm *jobs.Table) {
				jm.Add(1, "cat /path/to/fifo &")
				jm.MarkDone(1)
			},
			wantLines: []string{
				"[1]+  Done                    cat /path/to/fifo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			s := New(strings.NewReader(""), &out, io.Discard)
			tt.setup(s.state.Jobs)

			s.writeReapedJobs()

			gotLines := outputLines(out.String())
			if diff := cmp.Diff(tt.wantLines, gotLines, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("writeReapedJobs() output mismatch (-want +got):\n%s", diff)
			}

			if done := s.state.Jobs.ReapDone(); len(done) != 0 {
				t.Errorf("writeReapedJobs() left %d done jobs in table", len(done))
			}
		})
	}
}

func TestCommandFound(t *testing.T) {
	tests := []struct {
		name       string
		fields     []string
		wantOK     bool
		wantNotFnd string
	}{
		{name: "empty fields", fields: nil},
		{name: "builtin", fields: []string{"echo", "hello"}, wantOK: true},
		{name: "missing", fields: []string{"missing"}, wantNotFnd: "missing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNotFound, gotOK := commandFound(tt.fields)
			if gotOK != tt.wantOK {
				t.Errorf("commandFound() ok = %v, want %v", gotOK, tt.wantOK)
			}
			if diff := cmp.Diff(tt.wantNotFnd, gotNotFound); diff != "" {
				t.Errorf("commandFound() notFound mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestValidatePipelineSegments(t *testing.T) {
	tests := []struct {
		name       string
		segments   [][]string
		wantOK     bool
		wantNotFnd string
	}{
		{
			name: "all builtins",
			segments: [][]string{
				{"echo", "hello"},
				{"type", "echo"},
			},
			wantOK: true,
		},
		{
			name: "missing command in pipeline",
			segments: [][]string{
				{"echo", "hello"},
				{"missing"},
			},
			wantNotFnd: "missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNotFound, gotOK := validatePipelineSegments(tt.segments)
			if gotOK != tt.wantOK {
				t.Errorf("validatePipelineSegments() ok = %v, want %v", gotOK, tt.wantOK)
			}
			if diff := cmp.Diff(tt.wantNotFnd, gotNotFound); diff != "" {
				t.Errorf("validatePipelineSegments() notFound mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestExecuteLine(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		wantStop   bool
		wantStdout []string
		wantErr    bool
	}{
		{
			name:     "empty line",
			line:     "",
			wantStop: false,
		},
		{
			name:     "exit stops shell",
			line:     "exit",
			wantStop: true,
		},
		{
			name: "echo prints arguments",
			line: "echo hello",
			wantStdout: []string{
				"hello",
			},
		},
		{
			name: "history lists previous commands",
			line: "history",
			wantStdout: []string{
				"    1  history",
			},
		},
		{
			name: "history limit shows last entry",
			line: "history 1",
			wantStdout: []string{
				"    1  history 1",
			},
		},
		{
			name: "unknown command prints error",
			line: "missing",
			wantStdout: []string{
				"missing: command not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			s := New(strings.NewReader(""), &out, io.Discard)

			gotStop, err := s.ExecuteLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ExecuteLine(%q) error = %v, wantErr %v", tt.line, err, tt.wantErr)
			}
			if gotStop != tt.wantStop {
				t.Errorf("ExecuteLine(%q) stop = %v, want %v", tt.line, gotStop, tt.wantStop)
			}

			gotStdout := outputLines(out.String())
			if diff := cmp.Diff(tt.wantStdout, gotStdout, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ExecuteLine(%q) stdout mismatch (-want +got):\n%s", tt.line, diff)
			}
		})
	}
}

func outputLines(text string) []string {
	text = strings.TrimSuffix(text, "\n")
	if text == "" {
		return nil
	}
	return strings.Split(text, "\n")
}
