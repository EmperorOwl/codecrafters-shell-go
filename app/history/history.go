package history

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/codecrafters-io/shell-starter-go/app/files"
)

// Entry is a numbered command from the shell history.
type Entry struct {
	Number  int
	Command string
}

// HistoryList stores executed commands for the history builtin.
type HistoryList struct {
	mu       sync.Mutex
	commands []string
}

// Add records a command line in history.
func (l *HistoryList) Add(command string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.commands = append(l.commands, command)
}

// List returns a snapshot of all history entries with line numbers.
func (l *HistoryList) List() []Entry {
	return l.listEntries(0)
}

// ListLast returns the last n history entries, preserving original line numbers.
// If n is zero or greater than the number of entries, all entries are returned.
func (l *HistoryList) ListLast(n int) []Entry {
	return l.listEntries(n)
}

// Previous returns the command stepsBack entries before the most recent one.
// stepsBack 0 is the most recent command.
func (l *HistoryList) Previous(stepsBack int) (string, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if stepsBack < 0 || stepsBack >= len(l.commands) {
		return "", false
	}
	return l.commands[len(l.commands)-1-stepsBack], true
}

// ReadFromFile appends commands from the given file to the history list.
// Empty lines are skipped.
func (l *HistoryList) ReadFromFile(path string) error {
	lines, err := files.ReadLines(path)
	if err != nil {
		return err
	}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		l.Add(line)
	}
	return nil
}

// WriteToFile writes all commands in the history list to path.
func (l *HistoryList) WriteToFile(path string) error {
	l.mu.Lock()
	commands := append([]string(nil), l.commands...)
	l.mu.Unlock()
	return files.WriteLines(path, commands)
}

func (l *HistoryList) listEntries(limit int) []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()

	start := 0
	if limit > 0 && limit < len(l.commands) {
		start = len(l.commands) - limit
	}

	entries := make([]Entry, len(l.commands)-start)
	for i, command := range l.commands[start:] {
		entries[i] = Entry{
			Number:  start + i + 1,
			Command: command,
		}
	}
	return entries
}

// FormatLines returns bash-style display lines for the given entries.
func FormatLines(entries []Entry) []string {
	lines := make([]string, len(entries))
	for i, entry := range entries {
		lines[i] = formatLine(entry)
	}
	return lines
}

// WriteAll prints each history entry on its own line using bash-style formatting.
func WriteAll(out io.Writer, entries []Entry) {
	for _, line := range FormatLines(entries) {
		fmt.Fprintln(out, line)
	}
}

func formatLine(entry Entry) string {
	return fmt.Sprintf("%5d  %s", entry.Number, entry.Command)
}
