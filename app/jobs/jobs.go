package jobs

import (
	"fmt"
	"io"
	"strings"
)

type Job struct {
	Number  int
	PID     int
	Command string
	Status  string
}

func AddJob(jobs *[]Job, nextID *int, pid int, command string) int {
	*nextID++
	job := Job{
		Number:  *nextID,
		PID:     pid,
		Command: command,
		Status:  "Running",
	}
	*jobs = append(*jobs, job)
	return job.Number
}

func WriteAll(out io.Writer, jobs []Job) {
	for i, job := range jobs {
		fmt.Fprintln(out, formatLine(job, i, len(jobs)))
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
