package executor

import (
	"bytes"
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/session"
)

func (e *Executor) runPipeline(segments [][]string, stdout, stderr io.Writer, state *session.State) error {
	n := len(segments)
	readers := make([]io.ReadCloser, n-1)
	writers := make([]io.WriteCloser, n-1)
	for i := 0; i < n-1; i++ {
		readers[i], writers[i] = io.Pipe()
	}

	results := make(chan error, n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer func() {
				if i < n-1 {
					_ = writers[i].Close()
				}
			}()

			out := stdout
			if i < len(writers) {
				out = writers[i]
			}

			fields := segments[i]
			stdin := pipelineStdin(i, fields, readers)

			var err error
			if builtins.IsBuiltin(fields[0]) {
				_, err = e.runBuiltin(out, stderr, state, fields, stdin)
			} else {
				err = e.runExternal(out, stderr, fields, stdin)
			}

			results <- err
		}()
	}

	var lastErr error
	for i := 0; i < n; i++ {
		err := <-results
		if i == n-1 {
			lastErr = err
		}
	}
	return lastErr
}

func pipelineStdin(stage int, fields []string, readers []io.ReadCloser) io.Reader {
	if stage > 0 {
		return readers[stage-1]
	}
	if !builtins.IsBuiltin(fields[0]) {
		return bytes.NewReader(nil)
	}
	return nil
}
