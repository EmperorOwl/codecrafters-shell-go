package shell

import (
	"bytes"
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
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

	pipelineCommands, notFound, ok := s.buildPipelineCommands(commands, outputs.Stdout, outputs.Stderr, pipeWriters)
	if !ok {
		if notFound != "" {
			ctx.printCommandNotFound(notFound)
		}
		return false, nil
	}

	if err := nonExitError(s.runPipelineCommands(pipelineCommands, readers, writers)); err != nil {
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

func (s *Shell) buildPipelineCommands(segments [][]string, stdout, stderr io.Writer, writers []io.Writer) ([]resolvedCommand, string, bool) {
	commands := make([]resolvedCommand, len(segments))
	for i, fields := range segments {
		if len(fields) == 0 {
			return nil, "", false
		}

		out := stdout
		if i < len(writers) {
			out = writers[i]
		}

		cmd, notFound, ok := s.resolveCommand(fields, out, stderr)
		if !ok {
			return nil, notFound, false
		}
		commands[i] = cmd
	}
	return commands, "", true
}

func (s *Shell) runPipelineCommands(commands []resolvedCommand, readers []io.ReadCloser, writers []io.WriteCloser) error {
	n := len(commands)
	results := make(chan error, n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer func() {
				if i < n-1 {
					_ = writers[i].Close()
				}
			}()

			var err error
			switch {
			case commands[i].builtin != nil:
				if i > 0 {
					_, err = runDrainingStdin(commands[i].builtin, readers[i-1])
				} else {
					_, err = commands[i].builtin.Run()
				}
			case commands[i].external != nil:
				if i > 0 {
					commands[i].external.Stdin = readers[i-1]
				} else {
					commands[i].external.Stdin = bytes.NewReader(nil)
				}
				err = commands[i].external.Run()
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

func runDrainingStdin(cmd *builtins.Builtin, stdin io.Reader) (bool, error) {
	drainDone := make(chan struct{})
	go func() {
		_, _ = io.Copy(io.Discard, stdin)
		close(drainDone)
	}()

	exitShell, err := cmd.Run()
	<-drainDone
	return exitShell, err
}
