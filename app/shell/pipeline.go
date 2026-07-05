package shell

import (
	"bytes"
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/external"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

func (s *Shell) executePipeline(segments [][]string, ctx lineContext) (bool, error) {
	if len(segments) < 2 {
		return false, nil
	}

	commands, redirect := parsePipelineSegments(segments)
	outputs, err := openCommandOutputs(ctx.stdout, ctx.stderr, redirect)
	if err != nil {
		return true, err
	}
	defer outputs.Close()

	n := len(commands)
	readers := make([]io.ReadCloser, n-1)
	writers := make([]io.WriteCloser, n-1)
	for i := 0; i < n-1; i++ {
		readers[i], writers[i] = io.Pipe()
	}

	pipeWriters := make([]io.Writer, n-1)
	for i := range writers {
		pipeWriters[i] = writers[i]
	}

	notFound, ok := validatePipelineSegments(commands)
	if !ok {
		if notFound != "" {
			ctx.printCommandNotFound(notFound)
		}
		return false, nil
	}

	if err := nonExitError(s.runPipelineCommands(commands, outputs.Stdout, outputs.Stderr, pipeWriters, readers, writers)); err != nil {
		return true, err
	}
	return false, nil
}

func parsePipelineSegments(segments [][]string) ([][]string, parser.Redirect) {
	commands := make([][]string, len(segments))
	var redirect parser.Redirect
	for i, segment := range segments {
		fields, segmentRedirect := parser.ParseRedirect(segment)
		fields, _ = parser.StripBackground(fields)
		commands[i] = fields
		if i == len(segments)-1 {
			redirect = segmentRedirect
		}
	}
	return commands, redirect
}

func validatePipelineSegments(segments [][]string) (notFound string, ok bool) {
	for _, fields := range segments {
		notFound, ok := commandFound(fields)
		if !ok {
			return notFound, false
		}
	}
	return "", true
}

func (s *Shell) runPipelineCommands(segments [][]string, stdout, stderr io.Writer, writers []io.Writer, readers []io.ReadCloser, pipeWriters []io.WriteCloser) error {
	n := len(segments)
	results := make(chan error, n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer func() {
				if i < n-1 {
					_ = pipeWriters[i].Close()
				}
			}()

			out := stdout
			if i < len(writers) {
				out = writers[i]
			}

			fields := segments[i]
			var err error
			if builtins.IsBuiltin(fields[0]) {
				ctx := s.builtinContext(out, stderr)
				if i > 0 {
					_, err = runDrainingStdin(fields[0], fields[1:], ctx, readers[i-1])
				} else {
					_, err = builtins.Run(fields[0], fields[1:], ctx)
				}
			} else {
				prog, _ := external.New(fields, out, stderr)
				if i > 0 {
					prog.Stdin = readers[i-1]
				} else {
					prog.Stdin = bytes.NewReader(nil)
				}
				err = prog.Run()
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

func runDrainingStdin(name string, args []string, ctx *builtins.Context, stdin io.Reader) (bool, error) {
	drainDone := make(chan struct{})
	go func() {
		_, _ = io.Copy(io.Discard, stdin)
		close(drainDone)
	}()

	exitShell, err := builtins.Run(name, args, ctx)
	<-drainDone
	return exitShell, err
}
