package completer

import (
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/completion"
	"github.com/codecrafters-io/shell-starter-go/app/repl"
	"github.com/codecrafters-io/shell-starter-go/app/terminal"
)

// Completer orchestrates tab completion candidate sourcing and bash-style tab behavior.
type Completer struct {
	state *repl.State
}

// New returns a completer wired to the given REPL state.
func New(state *repl.State) *Completer {
	return &Completer{state: state}
}

// HandleTab computes a completion for the current buffer and applies bash-style tab
// behavior (extend buffer, ring bell, or show listings on a second tab).
func (c *Completer) HandleTab(state *terminal.TabState, buffer string) terminal.TabResult {
	newBuffer, listings := c.completeBuffer(buffer)
	return ApplyTabAction(state, buffer, newBuffer, listings)
}

// ApplyTabAction turns a completion result into terminal instructions.
// When listings are returned, the first tab rings the bell and stores them;
// a second tab with the same listings prints them. When the buffer changes,
// the new buffer replaces the input line.
func ApplyTabAction(state *terminal.TabState, buffer string, newBuffer string, listings []string) terminal.TabResult {
	if len(listings) > 0 {
		if slices.Equal(state.PendingListings, listings) {
			state.PendingListings = nil
			return terminal.TabResult{Buffer: buffer, ListingsToShow: listings}
		}
		state.PendingListings = listings
		return terminal.TabResult{Buffer: buffer, RingBell: true}
	}

	state.PendingListings = nil
	if newBuffer != buffer {
		return terminal.TabResult{Buffer: newBuffer}
	}
	return terminal.TabResult{Buffer: buffer, RingBell: true}
}

func (c *Completer) completeBuffer(buffer string) (newBuffer string, listings []string) {
	if !strings.Contains(buffer, " ") {
		return c.completeCommand(buffer)
	}

	if candidates, ok := c.programmableCandidates(buffer); ok {
		return c.completeArgument(buffer, candidates)
	}

	return c.completeFile(buffer)
}

func (c *Completer) programmableCandidates(buffer string) ([]string, bool) {
	opts := BuildCompleterOptions(buffer)
	if c.state == nil || c.state.Completion == nil {
		return nil, false
	}
	scriptPath, ok := c.state.Completion.Lookup(opts.Command)
	if !ok {
		return nil, false
	}
	opts.Path = scriptPath

	candidates, err := completion.RunCompleter(opts)
	if err != nil {
		return []string{}, true
	}
	return candidates, true
}

// BuildCompleterOptions parses the input buffer into the arguments and
// environment variables expected by a programmable completer script.
func BuildCompleterOptions(buffer string) completion.CompleterOptions {
	commandEnd := strings.Index(buffer, " ")
	if commandEnd < 0 {
		return completion.CompleterOptions{
			Command:   buffer,
			CompLine:  buffer,
			CompPoint: len(buffer),
		}
	}

	command := buffer[:commandEnd]
	afterCommand := buffer[commandEnd+1:]

	lastSpace := strings.LastIndex(afterCommand, " ")
	if lastSpace < 0 {
		currentWord := afterCommand
		previousWord := ""
		if currentWord != "" {
			previousWord = command
		}
		return completion.CompleterOptions{
			Command:      command,
			CurrentWord:  currentWord,
			PreviousWord: previousWord,
			CompLine:     buffer,
			CompPoint:    len(buffer),
		}
	}

	currentWord := afterCommand[lastSpace+1:]
	beforeCurrent := afterCommand[:lastSpace]
	prevLastSpace := strings.LastIndex(beforeCurrent, " ")
	previousWord := beforeCurrent
	if prevLastSpace >= 0 {
		previousWord = beforeCurrent[prevLastSpace+1:]
	}

	return completion.CompleterOptions{
		Command:      command,
		CurrentWord:  currentWord,
		PreviousWord: previousWord,
		CompLine:     buffer,
		CompPoint:    len(buffer),
	}
}
