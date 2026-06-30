package jobs

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAddJob(t *testing.T) {
	var jobList []Job
	nextID := 0

	number := AddJob(&jobList, &nextID, 42, "sleep 10 &")
	if number != 1 {
		t.Fatalf("AddJob() number = %d, want 1", number)
	}
	if nextID != 1 {
		t.Fatalf("nextID = %d, want 1", nextID)
	}
	want := []Job{{
		Number:  1,
		PID:     42,
		Command: "sleep 10 &",
		Status:  "Running",
	}}
	if diff := cmp.Diff(want, jobList); diff != "" {
		t.Errorf("job list mismatch (-want +got):\n%s", diff)
	}

	AddJob(&jobList, &nextID, 99, "sleep 5 &")
	if diff := cmp.Diff(2, jobList[1].Number); diff != "" {
		t.Errorf("second job number mismatch (-want +got):\n%s", diff)
	}
}

func TestWriteAll(t *testing.T) {
	tests := []struct {
		name      string
		jobs      []Job
		wantLines []string
	}{
		{
			name: "no jobs",
			jobs: nil,
		},
		{
			name: "one running job",
			jobs: []Job{{
				Number:  1,
				Command: "sleep 10 &",
				Status:  "Running",
			}},
			wantLines: []string{
				"[1]+  Running                 sleep 10 &",
			},
		},
		{
			name: "two jobs mark current and previous",
			jobs: []Job{
				{Number: 1, Command: "sleep 10 &", Status: "Running"},
				{Number: 2, Command: "sleep 20 &", Status: "Running"},
			},
			wantLines: []string{
				"[1]-  Running                 sleep 10 &",
				"[2]+  Running                 sleep 20 &",
			},
		},
		{
			name: "three jobs mark current, previous, and space",
			jobs: []Job{
				{Number: 1, Command: "sleep 10 &", Status: "Running"},
				{Number: 2, Command: "sleep 20 &", Status: "Running"},
				{Number: 3, Command: "sleep 30 &", Status: "Running"},
			},
			wantLines: []string{
				"[1]   Running                 sleep 10 &",
				"[2]-  Running                 sleep 20 &",
				"[3]+  Running                 sleep 30 &",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			jobList := tt.jobs
			WriteAllWithChecker(&out, &jobList, func(int) bool { return false })

			want := ""
			if len(tt.wantLines) > 0 {
				want = strings.Join(tt.wantLines, "\n") + "\n"
			}
			if diff := cmp.Diff(want, out.String()); diff != "" {
				t.Errorf("WriteAll() output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
