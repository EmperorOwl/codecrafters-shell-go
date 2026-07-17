package history

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/files"
	"github.com/codecrafters-io/shell-starter-go/app/testutils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestList_List(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*List)
		wantEntries []Entry
	}{
		{
			name: "empty history",
		},
		{
			name: "single command",
			setup: func(l *List) {
				l.Add("echo hello")
			},
			wantEntries: []Entry{{
				Number:  1,
				Command: "echo hello",
			}},
		},
		{
			name: "multiple commands",
			setup: func(l *List) {
				l.Add("previous_command_1")
				l.Add("previous_command_2")
				l.Add("history")
			},
			wantEntries: []Entry{
				{Number: 1, Command: "previous_command_1"},
				{Number: 2, Command: "previous_command_2"},
				{Number: 3, Command: "history"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := &List{}
			if tt.setup != nil {
				tt.setup(list)
			}

			got := list.List()
			if diff := cmp.Diff(tt.wantEntries, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("List() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestList_ListLast(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*List)
		limit int
		wantEntries []Entry
	}{
		{
			name:  "empty history",
			limit: 2,
		},
		{
			name: "limit greater than history length",
			setup: func(l *List) {
				l.Add("echo first")
			},
			limit: 2,
			wantEntries: []Entry{{
				Number:  1,
				Command: "echo first",
			}},
		},
		{
			name: "shows last two commands with original numbers",
			setup: func(l *List) {
				l.Add("echo first")
				l.Add("echo second")
				l.Add("history 2")
			},
			limit: 2,
			wantEntries: []Entry{
				{Number: 2, Command: "echo second"},
				{Number: 3, Command: "history 2"},
			},
		},
		{
			name: "zero limit returns all entries",
			setup: func(l *List) {
				l.Add("echo first")
				l.Add("echo second")
			},
			limit: 0,
			wantEntries: []Entry{
				{Number: 1, Command: "echo first"},
				{Number: 2, Command: "echo second"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := &List{}
			if tt.setup != nil {
				tt.setup(list)
			}

			got := list.ListLast(tt.limit)
			if diff := cmp.Diff(tt.wantEntries, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ListLast(%d) mismatch (-want +got):\n%s", tt.limit, diff)
			}
		})
	}
}

func TestList_Previous(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*List)
		stepsBack int
		want      string
		wantOK    bool
	}{
		{
			name:      "empty history",
			stepsBack: 0,
		},
		{
			name: "most recent command",
			setup: func(l *List) {
				l.Add("echo hello")
				l.Add("echo world")
			},
			stepsBack: 0,
			want:      "echo world",
			wantOK:    true,
		},
		{
			name: "earlier command",
			setup: func(l *List) {
				l.Add("echo hello")
				l.Add("echo world")
			},
			stepsBack: 1,
			want:      "echo hello",
			wantOK:    true,
		},
		{
			name: "before start of history",
			setup: func(l *List) {
				l.Add("echo hello")
			},
			stepsBack: 1,
		},
		{
			name:      "negative steps back",
			stepsBack: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := &List{}
			if tt.setup != nil {
				tt.setup(list)
			}

			got, ok := list.Previous(tt.stepsBack)
			if diff := cmp.Diff(tt.wantOK, ok); diff != "" {
				t.Errorf("Previous(%d) ok mismatch (-want +got):\n%s", tt.stepsBack, diff)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Previous(%d) mismatch (-want +got):\n%s", tt.stepsBack, diff)
			}
		})
	}
}

func TestList_AppendFromFile(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		missing     bool
		after       func(*List)
		wantEntries []Entry
	}{
		{
			name:        "appends commands from file",
			fileContent: "echo hello\necho world\n\n",
			after: func(l *List) {
				l.Add("history")
			},
			wantEntries: []Entry{
				{Number: 1, Command: "echo hello"},
				{Number: 2, Command: "echo world"},
				{Number: 3, Command: "history"},
			},
		},
		{
			name:    "missing file is ignored",
			missing: true,
		},
		{
			name:        "skips empty lines in file",
			fileContent: "echo hello\n\n  \necho world\n",
			wantEntries: []Entry{
				{Number: 1, Command: "echo hello"},
				{Number: 2, Command: "echo world"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := &List{}

			var path string
			if tt.missing {
				path = filepath.Join(t.TempDir(), "missing")
			} else {
				path = testutils.WriteTempFile(t, "histfile", tt.fileContent)
			}

			if err := list.AppendFromFile(path); err != nil {
				t.Fatalf("AppendFromFile() error = %v", err)
			}
			if tt.after != nil {
				tt.after(list)
			}

			got := list.List()
			if diff := cmp.Diff(tt.wantEntries, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("AppendFromFile() history mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestList_ReadFromFile(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		setup       func(*List, string)
		after       func(*List)
		wantEntries func(string) []Entry
	}{
		{
			name:        "appends commands from file",
			fileContent: "echo hello\necho world\n\n",
			setup: func(l *List, path string) {
				l.Add("history -r " + path)
			},
			after: func(l *List) {
				l.Add("history")
			},
			wantEntries: func(path string) []Entry {
				return []Entry{
					{Number: 1, Command: "history -r " + path},
					{Number: 2, Command: "echo hello"},
					{Number: 3, Command: "echo world"},
					{Number: 4, Command: "history"},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := testutils.WriteTempFile(t, "histfile", tt.fileContent)

			list := &List{}
			if tt.setup != nil {
				tt.setup(list, path)
			}

			if err := list.ReadFromFile(path); err != nil {
				t.Fatalf("ReadFromFile() error = %v", err)
			}
			if tt.after != nil {
				tt.after(list)
			}

			got := list.List()
			if diff := cmp.Diff(tt.wantEntries(path), got); diff != "" {
				t.Errorf("ReadFromFile() history mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestList_ReadFromFileAppendSync(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		setup       func(*List)
		wantAppend  []string
	}{
		{
			name:        "append only writes commands added after read",
			fileContent: "echo loaded\n",
			setup: func(l *List) {
				l.Add("echo new")
			},
			wantAppend: []string{"echo new"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := testutils.WriteTempFile(t, "histfile", tt.fileContent)
			list := &List{}

			if err := list.ReadFromFile(path); err != nil {
				t.Fatalf("ReadFromFile() error = %v", err)
			}
			if tt.setup != nil {
				tt.setup(list)
			}

			appendPath := filepath.Join(t.TempDir(), "append.txt")
			if err := list.AppendToFile(appendPath); err != nil {
				t.Fatalf("AppendToFile() error = %v", err)
			}

			got, err := files.ReadLines(appendPath)
			if err != nil {
				t.Fatalf("ReadLines() error = %v", err)
			}
			if diff := cmp.Diff(tt.wantAppend, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("AppendToFile() after ReadFromFile mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestList_WriteToFile(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*List, string)
		wantFile func(string) []string
	}{
		{
			name: "writes all commands",
			setup: func(l *List, path string) {
				l.Add("echo hello")
				l.Add("echo world")
				l.Add("history -w " + path)
			},
			wantFile: func(path string) []string {
				return []string{
					"echo hello",
					"echo world",
					"history -w " + path,
				}
			},
		},
		{
			name: "empty history writes empty file",
			setup: func(*List, string) {},
			wantFile: func(string) []string {
				return []string{""}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "histfile")

			list := &List{}
			if tt.setup != nil {
				tt.setup(list, path)
			}

			if err := list.WriteToFile(path); err != nil {
				t.Fatalf("WriteToFile() error = %v", err)
			}

			got, err := files.ReadLines(path)
			if err != nil {
				t.Fatalf("ReadLines() error = %v", err)
			}
			if diff := cmp.Diff(tt.wantFile(path), got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("WriteToFile() content mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestList_AppendToFile(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		setup       func(t *testing.T, l *List, path string)
		wantFile    func(string) []string
	}{
		{
			name:        "appends new commands to existing file",
			fileContent: "echo initial_command_1\necho initial_command_2\n\n",
			setup: func(_ *testing.T, l *List, path string) {
				l.Add("echo new_command")
				l.Add("history -a " + path)
			},
			wantFile: func(path string) []string {
				return []string{
					"echo initial_command_1",
					"echo initial_command_2",
					"echo new_command",
					"history -a " + path,
				}
			},
		},
		{
			name: "appends only new commands since last sync",
			setup: func(t *testing.T, l *List, path string) {
				l.Add("echo first")
				l.Add("history -a " + path)
				if err := l.AppendToFile(path); err != nil {
					t.Fatalf("first AppendToFile() error = %v", err)
				}
				l.Add("echo second")
				l.Add("history -a " + path)
			},
			wantFile: func(path string) []string {
				return []string{
					"echo first",
					"history -a " + path,
					"echo second",
					"history -a " + path,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path string
			if tt.fileContent != "" {
				path = testutils.WriteTempFile(t, "histfile", tt.fileContent)
			} else {
				path = filepath.Join(t.TempDir(), "histfile")
			}

			list := &List{}
			if tt.setup != nil {
				tt.setup(t, list, path)
			}

			if err := list.AppendToFile(path); err != nil {
				t.Fatalf("AppendToFile() error = %v", err)
			}

			got, err := files.ReadLines(path)
			if err != nil {
				t.Fatalf("ReadLines() error = %v", err)
			}
			if diff := cmp.Diff(tt.wantFile(path), got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("AppendToFile() content mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestWriteAll(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*List)
		want  string
	}{
		{
			name: "empty history",
		},
		{
			name: "numbered commands",
			setup: func(l *List) {
				l.Add("previous_command_1")
				l.Add("previous_command_2")
				l.Add("history")
			},
			want: testutils.WantStdout([]string{
				"    1  previous_command_1",
				"    2  previous_command_2",
				"    3  history",
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := &List{}
			if tt.setup != nil {
				tt.setup(list)
			}

			var out bytes.Buffer
			WriteAll(&out, list.List())

			if diff := cmp.Diff(tt.want, out.String()); diff != "" {
				t.Errorf("WriteAll() output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
