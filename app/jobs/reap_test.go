package jobs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestReapJobs(t *testing.T) {
	tests := []struct {
		name          string
		jobs          []Job
		exited        func(int) bool
		wantDisplay   []Job
		wantRemaining []Job
	}{
		{
			name: "exited job shown as done and removed",
			jobs: []Job{{Number: 1, PID: 1, Command: "sleep 1 &", Status: "Running"}},
			exited: func(int) bool {
				return true
			},
			wantDisplay: []Job{{
				Number:  1,
				PID:     1,
				Command: "sleep 1",
				Status:  "Done",
			}},
		},
		{
			name: "running job kept in table",
			jobs: []Job{{Number: 1, PID: 1, Command: "cat fifo &", Status: "Running"}},
			exited: func(int) bool {
				return false
			},
			wantDisplay: []Job{{
				Number:  1,
				PID:     1,
				Command: "cat fifo &",
				Status:  "Running",
			}},
			wantRemaining: []Job{{
				Number:  1,
				PID:     1,
				Command: "cat fifo &",
				Status:  "Running",
			}},
		},
		{
			name: "only exited jobs are removed",
			jobs: []Job{
				{Number: 1, PID: 1, Command: "sleep 1 &", Status: "Running"},
				{Number: 2, PID: 2, Command: "sleep 2 &", Status: "Running"},
			},
			exited: func(pid int) bool {
				return pid == 1
			},
			wantDisplay: []Job{
				{Number: 1, PID: 1, Command: "sleep 1", Status: "Done"},
				{Number: 2, PID: 2, Command: "sleep 2 &", Status: "Running"},
			},
			wantRemaining: []Job{
				{Number: 2, PID: 2, Command: "sleep 2 &", Status: "Running"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			display, remaining := reapJobs(tt.jobs, tt.exited)
			if diff := cmp.Diff(tt.wantDisplay, display, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("display mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantRemaining, remaining, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("remaining mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFormatDoneJobLine(t *testing.T) {
	job := Job{Number: 1, Command: "cat /path/to/fifo", Status: "Done"}
	want := "[1]+  Done                    cat /path/to/fifo"
	if diff := cmp.Diff(want, formatLine(job, 0, 1)); diff != "" {
		t.Errorf("formatLine() mismatch (-want +got):\n%s", diff)
	}
}
