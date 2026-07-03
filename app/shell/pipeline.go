package shell

import (
	"bytes"
	"io"

	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

type pipelineStage struct {
	fields    []string
	isBuiltin bool
}

type stageResult struct {
	err        error
	shouldExit bool
}

func (s *Shell) runBuiltin(fields []string, stdout, stderr io.Writer) (shouldExit bool) {
	_, shouldExit = TryBuiltin(fields, stdout, stderr, s.completers, &s.jobs)
	return shouldExit
}

func (s *Shell) runBuiltinDrainingStdin(fields []string, stdin io.Reader, stdout, stderr io.Writer) (shouldExit bool) {
	drainDone := make(chan struct{})
	go func() {
		_, _ = io.Copy(io.Discard, stdin)
		close(drainDone)
	}()

	shouldExit = s.runBuiltin(fields, stdout, stderr)
	<-drainDone
	return shouldExit
}

func resolvePipelineStages(segments [][]string) ([]pipelineStage, string, bool) {
	stages := make([]pipelineStage, len(segments))
	for i, fields := range segments {
		if len(fields) == 0 {
			return nil, "", false
		}
		if IsShellBuiltin(fields[0]) {
			stages[i] = pipelineStage{fields: fields, isBuiltin: true}
			continue
		}
		if _, ok := findExecutable(fields); ok {
			stages[i] = pipelineStage{fields: fields, isBuiltin: false}
			continue
		}
		return nil, fields[0], false
	}
	return stages, "", true
}

func (s *Shell) runPipelineStage(stage pipelineStage, stdin io.Reader, stdout, stderr io.Writer) (error, bool) {
	if stage.isBuiltin {
		if stdin != nil {
			return nil, s.runBuiltinDrainingStdin(stage.fields, stdin, stdout, stderr)
		}
		return nil, s.runBuiltin(stage.fields, stdout, stderr)
	}

	path, _ := findExecutable(stage.fields)
	cmd := newExternalCommand(stage.fields, path, stdout, stderr)
	if stdin != nil {
		cmd.Stdin = stdin
	} else {
		cmd.Stdin = bytes.NewReader(nil)
	}
	return cmd.Run(), false
}

func (s *Shell) runPipelineStages(stages []pipelineStage, stdout, stderr io.Writer) error {
	n := len(stages)
	readers := make([]io.ReadCloser, n-1)
	writers := make([]io.WriteCloser, n-1)
	for i := 0; i < n-1; i++ {
		readers[i], writers[i] = io.Pipe()
	}

	results := make(chan stageResult, n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			var in io.Reader
			if i > 0 {
				in = readers[i-1]
			}
			out := stdout
			if i < n-1 {
				out = writers[i]
			}
			defer func() {
				if i < n-1 {
					_ = writers[i].Close()
				}
			}()

			err, shouldExit := s.runPipelineStage(stages[i], in, out, stderr)
			results <- stageResult{err: err, shouldExit: shouldExit}
		}()
	}

	var lastErr error
	for i := 0; i < n; i++ {
		result := <-results
		if i == n-1 {
			lastErr = result.err
		}
	}
	return lastErr
}

func (s *Shell) ExecutePipeline(segments [][]string, stdout, stderr io.Writer) (executed bool, notFound string, err error) {
	if len(segments) < 2 {
		return false, "", nil
	}

	stages, notFound, ok := resolvePipelineStages(segments)
	if !ok {
		if notFound != "" {
			return false, notFound, nil
		}
		return false, "", nil
	}

	return true, "", s.runPipelineStages(stages, stdout, stderr)
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

func (s *Shell) executePipeline(segments [][]string, ctx lineContext) (bool, error) {
	commands, redirect := parsePipelineSegments(segments)
	outputs, err := openCommandOutputs(ctx.stdout, ctx.stderr, redirect)
	if err != nil {
		return true, err
	}
	defer outputs.Close()

	executed, notFound, execErr := s.ExecutePipeline(commands, outputs.Stdout, outputs.Stderr)
	if !executed {
		ctx.printCommandNotFound(notFound)
	} else if err := nonExitError(execErr); err != nil {
		return true, err
	}
	return ctx.stopAfter(nil)
}
