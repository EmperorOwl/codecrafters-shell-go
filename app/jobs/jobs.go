package jobs

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

const (
	StatusRunning = "Running"
	StatusDone    = "Done"
)

type Job struct {
	Number  int
	PID     int
	Command string
	Status  string
}

type JobTable struct {
	mu     sync.Mutex
	nextID int
	jobs   []Job
}

func (t *JobTable) Add(pid int, command string) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.nextID++
	job := Job{
		Number:  t.nextID,
		PID:     pid,
		Command: command,
		Status:  StatusRunning,
	}
	t.jobs = append(t.jobs, job)
	return job.Number
}

func (t *JobTable) MarkDone(jobNumber int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i, job := range t.jobs {
		if job.Number != jobNumber {
			continue
		}
		t.jobs[i].Status = StatusDone
		t.jobs[i].Command = strings.TrimSuffix(job.Command, " &")
		return
	}
}

func (t *JobTable) ListForDisplay() []Job {
	t.mu.Lock()
	defer t.mu.Unlock()

	display := make([]Job, len(t.jobs))
	copy(display, t.jobs)

	remaining := t.jobs[:0]
	for _, job := range t.jobs {
		if job.Status == StatusRunning {
			remaining = append(remaining, job)
		}
	}
	t.jobs = remaining
	return display
}

func WriteAll(out io.Writer, jobList []Job) {
	for i, job := range jobList {
		fmt.Fprintln(out, formatLine(job, i, len(jobList)))
	}
}

func markerForIndex(index, count int) string {
	switch {
	case index == count-1:
		return "+"
	case index == count-2:
		return "-"
	default:
		return " "
	}
}

func formatLine(job Job, index, count int) string {
	marker := markerForIndex(index, count)

	status := job.Status
	if len(status) < 24 {
		status += strings.Repeat(" ", 24-len(status))
	}

	return fmt.Sprintf("[%d]%s  %s%s", job.Number, marker, status, job.Command)
}
