package history

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestHistoryList_AddAndList(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		want     []Entry
	}{
		{
			name: "empty history",
		},
		{
			name:     "single command",
			commands: []string{"echo hello"},
			want: []Entry{{
				Number:  1,
				Command: "echo hello",
			}},
		},
		{
			name:     "multiple commands",
			commands: []string{"previous_command_1", "previous_command_2", "history"},
			want: []Entry{
				{Number: 1, Command: "previous_command_1"},
				{Number: 2, Command: "previous_command_2"},
				{Number: 3, Command: "history"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var list HistoryList
			for _, command := range tt.commands {
				list.Add(command)
			}

			got := list.List()
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("List() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHistoryList_ListLast(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		limit    int
		want     []Entry
	}{
		{
			name:  "empty history",
			limit: 2,
		},
		{
			name:     "limit greater than history length",
			commands: []string{"echo first"},
			limit:    2,
			want: []Entry{{
				Number:  1,
				Command: "echo first",
			}},
		},
		{
			name:     "shows last two commands with original numbers",
			commands: []string{"echo first", "echo second", "history 2"},
			limit:    2,
			want: []Entry{
				{Number: 2, Command: "echo second"},
				{Number: 3, Command: "history 2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var list HistoryList
			for _, command := range tt.commands {
				list.Add(command)
			}

			got := list.ListLast(tt.limit)
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ListLast(%d) mismatch (-want +got):\n%s", tt.limit, diff)
			}
		})
	}
}

func TestHistoryList_Previous(t *testing.T) {
	tests := []struct {
		name      string
		commands  []string
		stepsBack int
		want      string
		wantOK    bool
	}{
		{
			name:      "empty history",
			stepsBack: 0,
		},
		{
			name:      "most recent command",
			commands:  []string{"echo hello", "echo world"},
			stepsBack: 0,
			want:      "echo world",
			wantOK:    true,
		},
		{
			name:      "earlier command",
			commands:  []string{"echo hello", "echo world"},
			stepsBack: 1,
			want:      "echo hello",
			wantOK:    true,
		},
		{
			name:      "before start of history",
			commands:  []string{"echo hello"},
			stepsBack: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var list HistoryList
			for _, command := range tt.commands {
				list.Add(command)
			}

			got, ok := list.Previous(tt.stepsBack)
			if ok != tt.wantOK {
				t.Fatalf("Previous(%d) ok = %v, want %v", tt.stepsBack, ok, tt.wantOK)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Previous(%d) mismatch (-want +got):\n%s", tt.stepsBack, diff)
			}
		})
	}
}

func TestHistoryList_LoadHistfile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "histfile")
	content := "echo hello\necho world\n\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var list HistoryList
	if err := list.LoadHistfile(path); err != nil {
		t.Fatalf("LoadHistfile() error = %v", err)
	}
	list.Add("history")

	got := list.List()
	want := []Entry{
		{Number: 1, Command: "echo hello"},
		{Number: 2, Command: "echo world"},
		{Number: 3, Command: "history"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("LoadHistfile() history mismatch (-want +got):\n%s", diff)
	}
}

func TestHistoryList_LoadHistfileMissingFile(t *testing.T) {
	var list HistoryList
	if err := list.LoadHistfile(filepath.Join(t.TempDir(), "missing")); err != nil {
		t.Fatalf("LoadHistfile() error = %v", err)
	}
	if got := list.List(); len(got) != 0 {
		t.Fatalf("LoadHistfile() history = %v, want empty", got)
	}
}

func TestHistoryList_ReadFromFile(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/histfile"
	content := "echo hello\necho world\n\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var list HistoryList
	list.Add("history -r " + path)
	if err := list.ReadFromFile(path); err != nil {
		t.Fatalf("ReadFromFile() error = %v", err)
	}
	list.Add("history")

	got := list.List()
	want := []Entry{
		{Number: 1, Command: "history -r " + path},
		{Number: 2, Command: "echo hello"},
		{Number: 3, Command: "echo world"},
		{Number: 4, Command: "history"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ReadFromFile() history mismatch (-want +got):\n%s", diff)
	}
}

func TestHistoryList_WriteToFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "histfile")
	var list HistoryList
	list.Add("echo hello")
	list.Add("echo world")
	list.Add("history -w " + path)

	if err := list.WriteToFile(path); err != nil {
		t.Fatalf("WriteToFile() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	want := "echo hello\necho world\nhistory -w " + path + "\n"
	if diff := cmp.Diff(want, string(got)); diff != "" {
		t.Errorf("WriteToFile() content mismatch (-want +got):\n%s", diff)
	}
}

func TestHistoryList_AppendToFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "histfile")
	initial := "echo initial_command_1\necho initial_command_2\n\n"
	if err := os.WriteFile(path, []byte(initial), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var list HistoryList
	list.Add("echo new_command")
	list.Add("history -a " + path)

	if err := list.AppendToFile(path); err != nil {
		t.Fatalf("AppendToFile() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	want := strings.Join([]string{
		"echo initial_command_1",
		"echo initial_command_2",
		"echo new_command",
		"history -a " + path,
	}, "\n") + "\n"
	if diff := cmp.Diff(want, string(got)); diff != "" {
		t.Errorf("AppendToFile() content mismatch (-want +got):\n%s", diff)
	}
}

func TestHistoryList_AppendToFileOnlyNewCommands(t *testing.T) {
	path := filepath.Join(t.TempDir(), "histfile")

	var list HistoryList
	list.Add("echo first")
	list.Add("history -a " + path)
	if err := list.AppendToFile(path); err != nil {
		t.Fatalf("first AppendToFile() error = %v", err)
	}

	list.Add("echo second")
	list.Add("history -a " + path)
	if err := list.AppendToFile(path); err != nil {
		t.Fatalf("second AppendToFile() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	want := strings.Join([]string{
		"echo first",
		"history -a " + path,
		"echo second",
		"history -a " + path,
	}, "\n") + "\n"
	if diff := cmp.Diff(want, string(got)); diff != "" {
		t.Errorf("AppendToFile() content mismatch (-want +got):\n%s", diff)
	}
}

func TestWriteAll(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		want     string
	}{
		{
			name: "empty history",
		},
		{
			name:     "numbered commands",
			commands: []string{"previous_command_1", "previous_command_2", "history"},
			want: strings.Join([]string{
				"    1  previous_command_1",
				"    2  previous_command_2",
				"    3  history",
			}, "\n") + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var list HistoryList
			for _, command := range tt.commands {
				list.Add(command)
			}

			var out bytes.Buffer
			WriteAll(&out, list.List())

			if diff := cmp.Diff(tt.want, out.String()); diff != "" {
				t.Errorf("WriteAll() output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
