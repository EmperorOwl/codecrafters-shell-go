package parser

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "empty", input: "", want: nil},
		{name: "single word", input: "echo", want: []string{"echo"}},
		{
			name:  "spaces preserved in single quotes",
			input: "echo 'hello    world'",
			want:  []string{"echo", "hello    world"},
		},
		{
			name:  "unquoted spaces collapse between tokens",
			input: "echo hello    world",
			want:  []string{"echo", "hello", "world"},
		},
		{
			name:  "adjacent quoted strings concatenate",
			input: "echo 'hello''world'",
			want:  []string{"echo", "helloworld"},
		},
		{
			name:  "empty quotes are ignored",
			input: "echo hello''world",
			want:  []string{"echo", "helloworld"},
		},
		{
			name:  "quoted file paths for cat",
			input: "cat '/tmp/file name' '/tmp/file name with spaces'",
			want:  []string{"cat", "/tmp/file name", "/tmp/file name with spaces"},
		},
		{
			name:  "special characters literal in quotes",
			input: "echo '$HOME * ~'",
			want:  []string{"echo", "$HOME * ~"},
		},
		{
			name:  "spaces preserved in double quotes",
			input: `echo "hello    world"`,
			want:  []string{"echo", "hello    world"},
		},
		{
			name:  "adjacent double-quoted strings concatenate",
			input: `echo "hello""world"`,
			want:  []string{"echo", "helloworld"},
		},
		{
			name:  "quoted and unquoted strings concatenate",
			input: `echo "hello"world`,
			want:  []string{"echo", "helloworld"},
		},
		{
			name:  "separate double-quoted arguments",
			input: `echo "hello" "world"`,
			want:  []string{"echo", "hello", "world"},
		},
		{
			name:  "single quotes literal inside double quotes",
			input: `echo "shell's test"`,
			want:  []string{"echo", "shell's test"},
		},
		{
			name:  "quoted file paths with mixed quotes for cat",
			input: `cat "/tmp/file name" "/tmp/'file name' with spaces"`,
			want:  []string{"cat", "/tmp/file name", "/tmp/'file name' with spaces"},
		},
		{
			name:  "multiple double-quoted arguments with internal spaces",
			input: `echo "quz  hello"  "bar"`,
			want:  []string{"echo", "quz  hello", "bar"},
		},
		{
			name:  "escaped spaces form one argument",
			input: `echo three\ \ \ spaces`,
			want:  []string{"echo", "three   spaces"},
		},
		{
			name:  "escaped space then unescaped spaces split arguments",
			input: `echo before\     after`,
			want:  []string{"echo", "before ", "after"},
		},
		{
			name:  "backslash escapes regular letter",
			input: `echo test\nexample`,
			want:  []string{"echo", "testnexample"},
		},
		{
			name:  "backslash escapes backslash",
			input: `echo hello\\world`,
			want:  []string{"echo", `hello\world`},
		},
		{
			name:  "escaped single quotes outside quotes",
			input: `echo \'hello\'`,
			want:  []string{"echo", `'hello'`},
		},
		{
			name:  "escaped quotes split by unescaped space",
			input: `echo \'\"literal quotes\"\'`,
			want:  []string{"echo", `'"literal`, `quotes"'`},
		},
		{
			name:  "escaped digit in path",
			input: `cat /tmp/ignore_\2`,
			want:  []string{"cat", "/tmp/ignore_2"},
		},
		{
			name:  "escaped backslash before digit in path",
			input: `cat /tmp/just_one_\\3`,
			want:  []string{"cat", `/tmp/just_one_\3`},
		},
		{
			name:  "literal backslash before non-special char in double quotes",
			input: `echo "A \ escapes itself"`,
			want:  []string{"echo", `A \ escapes itself`},
		},
		{
			name:  "escaped double quote inside double quotes",
			input: `echo "A \" inside double quotes"`,
			want:  []string{"echo", `A " inside double quotes`},
		},
		{
			name:  "literal backslash-n inside double quotes",
			input: `echo "just'one'\n'backslash"`,
			want:  []string{"echo", `just'one'\n'backslash`},
		},
		{
			name:  "escaped quotes inside and outside quoted segments",
			input: `echo "inside\"literal_quote."outside\"`,
			want:  []string{"echo", `inside"literal_quote.outside"`},
		},
		{
			name:  "quoted path with space",
			input: `cat /tmp/"number 1"`,
			want:  []string{"cat", `/tmp/number 1`},
		},
		{
			name:  "quoted path with escaped double quote",
			input: `cat /tmp/"doublequote \" 2"`,
			want:  []string{"cat", `/tmp/doublequote " 2`},
		},
		{
			name:  "quoted path with literal backslash",
			input: `cat /tmp/"backslash \ 3"`,
			want:  []string{"cat", `/tmp/backslash \ 3`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Tokenize(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokenize(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
