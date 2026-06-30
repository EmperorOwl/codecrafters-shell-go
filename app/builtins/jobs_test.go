package builtins

import (
	"bytes"
	"strings"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/jobs"
	"github.com/google/go-cmp/cmp"
)

func TestJobs(t *testing.T) {
	neverExited := func(int) bool { return false }

	tests := []struct {
		name      string
		jobs      []jobs.Job
		hasExited func(int) bool
		wantLines []string
	}{
		{
			name:      "no jobs",
			jobs:      nil,
			hasExited: neverExited,
		},
		{
			name: "one running job",
			jobs: []jobs.Job{{
				Number:  1,
				Command: "sleep 10 &",
				Status:  "Running",
			}},
			hasExited: neverExited,
			wantLines: []string{
				"[1]+  Running                 sleep 10 &",
			},
		},
		{
			name: "reaped job shown as done and removed",
			jobs: []jobs.Job{{
				Number:  1,
				PID:     1,
				Command: "cat /path/to/fifo &",
				Status:  "Running",
			}},
			hasExited: func(int) bool { return true },
			wantLines: []string{
				"[1]+  Done                    cat /path/to/fifo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			jobList := tt.jobs
			Jobs(&out, &jobList, tt.hasExited)

			want := ""
			if len(tt.wantLines) > 0 {
				want = strings.Join(tt.wantLines, "\n") + "\n"
			}
			if diff := cmp.Diff(want, out.String()); diff != "" {
				t.Errorf("Jobs() output mismatch (-want +got):\n%s", diff)
			}
			if tt.name == "reaped job shown as done and removed" && len(jobList) != 0 {
				t.Errorf("Jobs() job list length = %d, want 0 after reaping", len(jobList))
			}
		})
	}
}
