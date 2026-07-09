package history

import (
	"fmt"
	"io"
	"sync"
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
	l.mu.Lock()
	defer l.mu.Unlock()

	entries := make([]Entry, len(l.commands))
	for i, command := range l.commands {
		entries[i] = Entry{
			Number:  i + 1,
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
