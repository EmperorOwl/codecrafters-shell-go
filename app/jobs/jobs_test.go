package jobs

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestJobTableAdd(t *testing.T) {
	var table JobTable

	number := table.Add(42, "sleep 10 &")
	if number != 1 {
		t.Fatalf("Add() number = %d, want 1", number)
	}

	table.mu.Lock()
	want := []Job{{
		Number:  1,
		PID:     42,
		Command: "sleep 10 &",
		Status:  StatusRunning,
	}}
	if diff := cmp.Diff(want, table.jobs); diff != "" {
		t.Errorf("jobs mismatch (-want +got):\n%s", diff)
	}
	table.mu.Unlock()

	table.Add(99, "sleep 5 &")
	table.mu.Lock()
	if table.jobs[1].Number != 2 {
		t.Errorf("second job number = %d, want 2", table.jobs[1].Number)
	}
	table.mu.Unlock()
}

func TestJobTableMarkDone(t *testing.T) {
	var table JobTable
	table.Add(42, "sleep 1 &")

	table.MarkDone(1)

	table.mu.Lock()
	want := []Job{{
		Number:  1,
		PID:     42,
		Command: "sleep 1",
		Status:  StatusDone,
	}}
	if diff := cmp.Diff(want, table.jobs); diff != "" {
		t.Errorf("jobs mismatch (-want +got):\n%s", diff)
	}
	table.mu.Unlock()
}

func TestJobTableListForDisplay(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*JobTable)
		wantDisplay   []Job
		wantRemaining []Job
	}{
		{
			name: "running job stays in table",
			setup: func(t *JobTable) {
				t.Add(1, "cat fifo &")
			},
			wantDisplay: []Job{{
				Number:  1,
				PID:     1,
				Command: "cat fifo &",
				Status:  StatusRunning,
			}},
			wantRemaining: []Job{{
				Number:  1,
				PID:     1,
				Command: "cat fifo &",
				Status:  StatusRunning,
			}},
		},
		{
			name: "done job shown once then removed",
			setup: func(t *JobTable) {
				t.Add(1, "sleep 1 &")
				t.MarkDone(1)
			},
			wantDisplay: []Job{{
				Number:  1,
				PID:     1,
				Command: "sleep 1",
				Status:  StatusDone,
			}},
		},
		{
			name: "only done jobs are removed",
			setup: func(t *JobTable) {
				t.Add(1, "sleep 1 &")
				t.Add(2, "sleep 2 &")
				t.MarkDone(1)
			},
			wantDisplay: []Job{
				{Number: 1, PID: 1, Command: "sleep 1", Status: StatusDone},
				{Number: 2, PID: 2, Command: "sleep 2 &", Status: StatusRunning},
			},
			wantRemaining: []Job{
				{Number: 2, PID: 2, Command: "sleep 2 &", Status: StatusRunning},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var table JobTable
			tt.setup(&table)

			display := table.ListForDisplay()
			if diff := cmp.Diff(tt.wantDisplay, display, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("display mismatch (-want +got):\n%s", diff)
			}

			table.mu.Lock()
			remaining := append([]Job(nil), table.jobs...)
			table.mu.Unlock()
			if diff := cmp.Diff(tt.wantRemaining, remaining, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("remaining mismatch (-want +got):\n%s", diff)
			}
		})
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
				Status:  StatusRunning,
			}},
			wantLines: []string{
				"[1]+  Running                 sleep 10 &",
			},
		},
		{
			name: "two jobs mark current and previous",
			jobs: []Job{
				{Number: 1, Command: "sleep 10 &", Status: StatusRunning},
				{Number: 2, Command: "sleep 20 &", Status: StatusRunning},
			},
			wantLines: []string{
				"[1]-  Running                 sleep 10 &",
				"[2]+  Running                 sleep 20 &",
			},
		},
		{
			name: "three jobs mark current, previous, and space",
			jobs: []Job{
				{Number: 1, Command: "sleep 10 &", Status: StatusRunning},
				{Number: 2, Command: "sleep 20 &", Status: StatusRunning},
				{Number: 3, Command: "sleep 30 &", Status: StatusRunning},
			},
			wantLines: []string{
				"[1]   Running                 sleep 10 &",
				"[2]-  Running                 sleep 20 &",
				"[3]+  Running                 sleep 30 &",
			},
		},
		{
			name: "done job",
			jobs: []Job{{
				Number:  1,
				Command: "cat /path/to/fifo",
				Status:  StatusDone,
			}},
			wantLines: []string{
				"[1]+  Done                    cat /path/to/fifo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			WriteAll(&out, tt.jobs)

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
