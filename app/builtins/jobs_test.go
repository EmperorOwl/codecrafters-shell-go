package builtins

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestJobs(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*jobs.Table)
		wantLines []string
		wantJobs  []jobs.Job
	}{
		{
			name: "no jobs",
		},
		{
			name: "one running job",
			setup: func(t *jobs.Table) {
				t.Add(1, "sleep 10 &")
			},
			wantLines: []string{
				"[1]+  Running                 sleep 10 &",
			},
			wantJobs: []jobs.Job{{
				Number:  1,
				PID:     1,
				Command: "sleep 10 &",
				Status:  jobs.StatusRunning,
			}},
		},
		{
			name: "done job is reaped",
			setup: func(t *jobs.Table) {
				t.Add(1, "cat /path/to/fifo &")
				t.MarkDone(1)
			},
			wantLines: []string{
				"[1]+  Done                    cat /path/to/fifo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := jobs.NewTable()
			if tt.setup != nil {
				tt.setup(table)
			}

			var stdout bytes.Buffer
			jobsBuiltin(&stdout, table)

			if diff := cmp.Diff(wantStdout(tt.wantLines), stdout.String()); diff != "" {
				t.Errorf("jobsBuiltin() stdout mismatch (-want +got):\n%s", diff)
			}

			if tt.wantJobs != nil {
				if diff := cmp.Diff(tt.wantJobs, table.List(), cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("jobsBuiltin() remaining jobs mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
