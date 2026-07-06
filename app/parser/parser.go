package parser

// Command is a parsed single command with redirect and background markers applied.
type Command struct {
	Fields     []string
	Redirect   Redirect
	Background bool
}

// Line is a fully parsed input line: either one command or a pipeline.
type Line struct {
	Pipeline   bool
	Commands   [][]string
	Redirect   Redirect
	Background bool
}

// ParseCommand parses redirect and background markers from tokenized arguments.
func ParseCommand(tokens []string) Command {
	fields, redirect := ParseRedirect(tokens)
	fields, background := StripBackground(fields)
	return Command{
		Fields:     fields,
		Redirect:   redirect,
		Background: background,
	}
}

// ParsePipelineSegments parses redirect and background markers from pipeline segments.
// Redirects on the final segment apply to the whole pipeline; background markers are stripped.
func ParsePipelineSegments(segments [][]string) (commands [][]string, redirect Redirect) {
	commands = make([][]string, len(segments))
	for i, segment := range segments {
		fields, segmentRedirect := ParseRedirect(segment)
		fields, _ = StripBackground(fields)
		commands[i] = fields
		if i == len(segments)-1 {
			redirect = segmentRedirect
		}
	}
	return commands, redirect
}

// ParseLine tokenizes and parses a full input line into a command or pipeline.
func ParseLine(line string) Line {
	tokens := Tokenize(line)
	if segments := SplitPipelineTokens(tokens); len(segments) >= 2 {
		commands, redirect := ParsePipelineSegments(segments)
		return Line{
			Pipeline: true,
			Commands: commands,
			Redirect: redirect,
		}
	}

	cmd := ParseCommand(tokens)
	return Line{
		Commands:   [][]string{cmd.Fields},
		Redirect:   cmd.Redirect,
		Background: cmd.Background,
	}
}
