package executor

import (
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/jobs"
)

// Executor runs parsed commands using injected shell state and default I/O streams.
type Executor struct {
	jobTable           *jobs.JobTable
	completionRegistry *completion.CompletionRegistry
	stdin              io.Reader
	stdout             io.Writer
	stderr             io.Writer
}

// New returns an executor wired to the given shell state.
// Call SetIO with the terminal streams before executing commands.
func New(jobTable *jobs.JobTable, completionRegistry *completion.CompletionRegistry) *Executor {
	return &Executor{
		jobTable:           jobTable,
		completionRegistry: completionRegistry,
	}
}

// SetIO configures the default stdin, stdout, and stderr used for command execution.
func (e *Executor) SetIO(stdin io.Reader, stdout, stderr io.Writer) {
	e.stdin = stdin
	e.stdout = stdout
	e.stderr = stderr
}
