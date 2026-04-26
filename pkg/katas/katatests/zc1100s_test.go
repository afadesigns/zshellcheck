// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1100(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid basename with suffix flag",
			input:    `basename -s .txt file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid basename with multiple args",
			input:    `basename /path/to/file .txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid simple dirname",
			input: `dirname /path/to/file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1100",
					Message: "Use `${var%/*}` instead of `dirname` to extract the directory path. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid simple basename",
			input: `basename /path/to/file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1100",
					Message: "Use `${var##*/}` instead of `basename` to extract the filename. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1100")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1101(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid bc with file",
			input:    `bc script.bc`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid bc in pipeline",
			input: `bc -l`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1101",
					Message: "Use `$(( ))` for arithmetic instead of `bc`. Zsh arithmetic expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid other command",
			input:    `calc 1+1`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1101")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1102(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "sudo redirection",
			input: `sudo echo "foo" > /etc/bar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1102",
					Message: "Redirection happens before `sudo`. This will likely fail permission checks. Use `| sudo tee`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "sudo append redirection",
			input: `sudo echo "foo" >> /etc/bar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1102",
					Message: "Redirection happens before `sudo`. This will likely fail permission checks. Use `| sudo tee`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid sudo usage",
			input:    `echo "foo" | sudo tee /etc/bar`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1102")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1103(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "PATH assignment",
			input: `PATH=$PATH:/usr/local/bin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1103",
					Message: "Use the `path` array instead of manually manipulating the `$PATH` string.",
					Line:    1,
					Column:  5,
				},
			},
		},
		/*
			{
				name:  "export PATH",
				input: `export PATH=$PATH:/bin`,
				expected: []katas.Violation{
					{
						KataID:  "ZC1103",
						Message: "Use the `path` array instead of manually manipulating the `$PATH` string.",
						Line:    1,
						Column:  1,
					},
				},
			},
		*/{
			name:     "path array assignment",
			input:    `path+=('/usr/local/bin')`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1103")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1104(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "export PATH assignment",
			input: `export PATH=$PATH:/bin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1104",
					Message: "Use the `path` array instead of manually manipulating the `$PATH` string.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid path array assignment",
			input:    `path+=('/usr/local/bin')`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1104")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1105(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "nested arithmetic expansion",
			input: `(( result = $((1+1)) + 2 ))`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1105",
					Message: "Avoid nested arithmetic expansions. Use intermediate variables for clarity.",
					Line:    1,
					Column:  2,
				},
			},
		},
		{
			name:     "simple arithmetic expansion",
			input:    `result=$((1+1))`,
			expected: []katas.Violation{},
		},
		{
			name:     "double parenthesis arithmetic command",
			input:    `(( 1 + 1 ))`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1105")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1106(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "set -x usage",
			input: `set -x`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1106",
					Message: "Avoid `set -x` in production scripts to prevent sensitive data exposure.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "set -eux usage",
			input: `set -eux`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1106",
					Message: "Avoid `set -x` in production scripts to prevent sensitive data exposure.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid set usage",
			input:    `set +x`,
			expected: []katas.Violation{},
		},
		{
			name:     "other set options",
			input:    `set -e`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1106")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1107(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want int
	}{
		{
			name: "Valid arithmetic",
			src:  "if (( a > b )); then echo yes; fi",
			want: 0,
		},
		{
			name: "String comparison",
			src:  "if [[ $a == $b ]]; then echo yes; fi",
			want: 0,
		},
		{
			name: "File check",
			src:  "if [[ -f file ]]; then echo yes; fi",
			want: 0,
		},
		{
			name: "Invalid -eq",
			src:  "if [ $a -eq $b ]; then echo yes; fi",
			want: 1,
		},
		{
			name: "Invalid -gt in [[ ]]",
			src:  "if [[ $a -gt 5 ]]; then echo yes; fi",
			want: 1,
		},
		{
			name: "Invalid -le",
			src:  "while [ $count -le 10 ]; do ((count++)); done",
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.src, "ZC1107")
			if len(violations) != tt.want {
				t.Errorf("Test %q failed: want %d violations, got %d", tt.name, tt.want, len(violations))
			}
		})
	}
}

func TestZC1108(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid tr delete",
			input:    `tr -d '\n'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid tr squeeze",
			input:    `tr -s ' '`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid tr with different sets",
			input:    `tr ':' '\n'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid tr lowercase to uppercase POSIX",
			input: `tr '[:lower:]' '[:upper:]'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1108",
					Message: "Use `${(U)var}` for case conversion instead of `tr`. Zsh parameter expansion flags avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid tr uppercase to lowercase range",
			input: `tr 'A-Z' 'a-z'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1108",
					Message: "Use `${(L)var}` for case conversion instead of `tr`. Zsh parameter expansion flags avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1108")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1109(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid cut with file",
			input:    `cut -d: -f1 /etc/passwd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cut with only field",
			input:    `cut -f1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cut in pipeline",
			input: `cut -d: -f1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1109",
					Message: "Use Zsh parameter expansion for field extraction instead of `cut`. `${var%%delim*}` or `${(s.delim.)var}` avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid cut with delimiter and field",
			input: `cut -d',' -f2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1109",
					Message: "Use Zsh parameter expansion for field extraction instead of `cut`. `${var%%delim*}` or `${(s.delim.)var}` avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1109")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1110(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid head with file",
			input:    `head -1 file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid head -5",
			input:    `head -5`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid tail with file",
			input:    `tail -1 file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid head -1 in pipeline",
			input: `head -1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1110",
					Message: "Use `${lines[1]}` instead of `head -1`. Zsh array subscripts avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid tail -1 in pipeline",
			input: `tail -1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1110",
					Message: "Use `${lines[-1]}` instead of `tail -1`. Zsh array subscripts avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid head -n 1",
			input: `head -n 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1110",
					Message: "Use `${lines[1]}` instead of `head -1`. Zsh array subscripts avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1110")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1111(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid xargs with null separator",
			input:    `xargs -0 rm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid xargs with parallel",
			input:    `xargs -P 4 cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid xargs with replace",
			input:    `xargs -I {} mv {} /dest`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid simple xargs",
			input: `xargs rm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1111",
					Message: "Consider using Zsh array iteration instead of `xargs`. `for item in ${(f)$(cmd)}` splits output by newlines without spawning xargs.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid xargs with command only",
			input: `xargs grep pattern`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1111",
					Message: "Consider using Zsh array iteration instead of `xargs`. `for item in ${(f)$(cmd)}` splits output by newlines without spawning xargs.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1111")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1112(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -c with file",
			input:    `grep -c pattern file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid grep without -c",
			input:    `grep pattern`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep -c in pipeline",
			input: `grep -c pattern`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1112",
					Message: "Use Zsh array filtering `${(M)array:#pattern}` or `${#${(f)...}}` for counting instead of `grep -c`. Avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid grep -v (not count)",
			input:    `grep -v pattern`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1112")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1113(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid readlink without resolve flag",
			input:    `readlink /path/to/link`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid realpath",
			input: `realpath /path/to/file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1113",
					Message: "Use `${var:A}` instead of `realpath` to resolve absolute paths. Zsh path modifiers avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid readlink -f",
			input: `readlink -f /path/to/link`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1113",
					Message: "Use `${var:A}` instead of `readlink -f` to resolve absolute paths. Zsh path modifiers avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid realpath with complex flags",
			input:    `realpath --relative-to /base /path`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1113")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1114(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid mktemp -d",
			input:    `mktemp -d`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid mktemp",
			input: `mktemp /tmp/zsh.XXXXXX`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1114",
					Message: "Consider using Zsh `=(cmd)` for temporary files instead of `mktemp`. Zsh auto-cleans temporary files created with `=(...)` process substitution.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid bare mktemp",
			input: `mktemp -t prefix`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1114",
					Message: "Consider using Zsh `=(cmd)` for temporary files instead of `mktemp`. Zsh auto-cleans temporary files created with `=(...)` process substitution.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1114")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1115(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid rev with file",
			input:    `rev file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid rev in pipeline",
			input: `rev -`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1115",
					Message: "Use Zsh string manipulation instead of `rev`. Parameter expansion can reverse strings without spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1115")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1116(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid tee with append",
			input:    `tee -a logfile`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid simple tee",
			input: `tee output.log`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1116",
					Message: "Use Zsh multios (`setopt multios`) instead of `tee`. With multios, `cmd > file1 > file2` writes to both files without spawning tee.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid tee with multiple files",
			input: `tee file1 file2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1116",
					Message: "Use Zsh multios (`setopt multios`) instead of `tee`. With multios, `cmd > file1 > file2` writes to both files without spawning tee.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1116")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1117(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid nohup usage",
			input: `nohup ./server`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1117",
					Message: "Use `cmd &!` or `cmd & disown` instead of `nohup cmd &`. Zsh `&!` is a built-in shorthand that avoids spawning nohup.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid nohup with redirect",
			input: `nohup ./server > /dev/null`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1117",
					Message: "Use `cmd &!` or `cmd & disown` instead of `nohup cmd &`. Zsh `&!` is a built-in shorthand that avoids spawning nohup.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1117")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1118(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid echo without -n",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid print -rn",
			input:    `print -rn hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid echo -n",
			input: `echo -n hello`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1118",
					Message: "Use `print -rn` instead of `echo -n`. `echo -n` behavior varies across shells; `print -rn` is the reliable Zsh idiom.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1118")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1119(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid date with format",
			input:    `date "+%Y-%m-%d"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid date no args",
			input:    `date -u`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid date +%s",
			input: `date '+%s'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1119",
					Message: "Use `$EPOCHSECONDS` or `$EPOCHREALTIME` (via `zmodload zsh/datetime`) instead of `date +%s`. Avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1119")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1120(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid pwd -P",
			input:    `pwd -P`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid bare pwd",
			input: `pwd -L`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1120",
					Message: "Use `$PWD` instead of `pwd`. Zsh maintains `$PWD` as a built-in variable, avoiding an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1120")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1121(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid hostname -f",
			input:    `hostname -f`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid hostname -s",
			input:    `hostname -s`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid simple hostname",
			input: `hostname myhost`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1121",
					Message: "Use `$HOST` instead of `hostname`. Zsh maintains `$HOST` as a built-in variable, avoiding an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1121")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1122(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid whoami",
			input: `whoami`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1122",
					Message: "Use `$USER` instead of `whoami`. Zsh maintains `$USER` as a built-in variable, avoiding an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid other identifier",
			input:    `mycommand`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1122")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1123(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid uname -r",
			input:    `uname -r`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid uname -m",
			input:    `uname -m`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid uname -s",
			input: `uname -s`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1123",
					Message: "Use `$OSTYPE` instead of `uname -s` for OS detection. Zsh maintains `$OSTYPE` as a built-in variable.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1123")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1124(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid cat with file",
			input:    `cat readme.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat /dev/null",
			input: `cat /dev/null`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1124",
					Message: "Use `: > file` instead of `cat /dev/null > file` to truncate. The `:` builtin avoids spawning cat.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1124")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1125(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -q with file",
			input:    `grep -q pattern file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid grep without -q",
			input:    `grep pattern`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep -q in pipeline",
			input: `grep -q pattern`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1125",
					Message: "Use `[[ $var =~ pattern ]]` or `[[ $var == *pattern* ]]` instead of piping through `grep -q`. Zsh pattern matching avoids spawning external processes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1125")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1126(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sort -u",
			input:    `sort -u file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sort | uniq -c",
			input:    `sort file | uniq -c`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sort | uniq",
			input: `sort file | uniq`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1126",
					Message: "Use `sort -u` instead of `sort | uniq`. Combining into one command avoids an unnecessary pipeline.",
					Line:    1,
					Column:  11,
				},
			},
		},
		{
			name:     "non-pipe operator",
			input:    `echo a && echo b`,
			expected: []katas.Violation{},
		},
		{
			name:     "pipe but not sort",
			input:    `cat file | uniq`,
			expected: []katas.Violation{},
		},
		{
			name:     "sort piped to non-uniq",
			input:    `sort file | head`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1126")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1127(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ls -la",
			input:    `ls -la`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ls -1",
			input: `ls -1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1127",
					Message: "Use Zsh glob qualifiers `files=(*(N)); echo ${#files}` instead of `ls -1 | wc -l`. Avoids spawning external processes for file counting.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1127")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1128(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid touch with timestamp",
			input:    `touch -t 202301011200 file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid touch multiple files",
			input:    `touch file1 file2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid touch single file",
			input: `touch newfile.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1128",
					Message: "Use `> file` instead of `touch file` to create an empty file. This avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1128")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1129(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid wc -l with file",
			input:    `wc -l file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid wc -c without file",
			input:    `wc -c`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid wc -c with file",
			input: `wc -c file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1129",
					Message: "Use `zstat +size file` (via `zmodload zsh/stat`) instead of `wc -c file`. Avoids reading the entire file for a simple size query.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1129")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1131(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid cat with grep",
			input:    `cat file | grep pattern`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat | read",
			input: `cat file.txt | read line`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1131",
					Message: "Use `while read line; do ...; done < file` instead of `cat file | while read line`. Avoids unnecessary cat and subshell from the pipe.",
					Line:    1,
					Column:  14,
				},
			},
		},
		{
			name:     "non-pipe operator",
			input:    `echo hello && echo world`,
			expected: []katas.Violation{},
		},
		{
			name:     "pipe but left not cat",
			input:    `echo hello | read line`,
			expected: []katas.Violation{},
		},
		{
			name:     "cat no args piped to read",
			input:    `cat | read line`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1131")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1132(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -o with file",
			input:    `grep -o pattern file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid grep without -o",
			input:    `grep pattern`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep -o in pipeline",
			input: `grep -o pattern`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1132",
					Message: "Use Zsh pattern extraction `${(M)var:#pattern}` or `[[ $var =~ regex ]] && echo $match[1]` instead of piping through `grep -o`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1132")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1133(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid kill SIGTERM",
			input:    `kill 1234`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid kill -15",
			input:    `kill -15 1234`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid kill -9",
			input: `kill -9 1234`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1133",
					Message: "Avoid `kill -9` as a first resort. Use `kill` (SIGTERM) first to allow graceful shutdown, then escalate to `kill -9` if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid kill -KILL",
			input: `kill -KILL 1234`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1133",
					Message: "Avoid `kill -9` as a first resort. Use `kill` (SIGTERM) first to allow graceful shutdown, then escalate to `kill -9` if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1133")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1134(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sleep 5",
			input:    `sleep 5`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sleep 30",
			input:    `sleep 30`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sleep 0.1",
			input: `sleep 0.1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1134",
					Message: "Avoid `sleep 0.1` in loops. Short sleep intervals suggest busy-waiting. Consider event-driven alternatives like `inotifywait` or `zle -F`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1134")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1135(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid env -i",
			input:    `env -i cmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid env VAR=val cmd",
			input: `env FOO=bar cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1135",
					Message: "Use inline `VAR=val cmd` instead of `env VAR=val cmd`. Zsh supports inline env assignment without spawning env.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1135")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1136(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid rm -rf with literal path",
			input:    `rm -rf /tmp/build`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid rm without -rf",
			input:    `rm $file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid rm -rf with variable",
			input: `rm -rf $dir`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1136",
					Message: "Avoid `rm -rf $var` without safeguards. Use `rm -rf ${var:?}` to abort if the variable is empty, preventing accidental deletion.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1136")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1137(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid mktemp",
			input:    `mktemp /tmp/foo.XXXXXX`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid temp with variable",
			input:    `cat $tmpfile`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid hardcoded tmp path",
			input: `cp data /tmp/myapp.log`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1137",
					Message: "Avoid hardcoded `/tmp/` paths. Use `mktemp` or Zsh `=(cmd)` for temp files to prevent race conditions and symlink attacks.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1137")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1139(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid source local file",
			input:    `source /usr/local/lib/utils.zsh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid source URL",
			input: `source https://example.com/script.zsh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1139",
					Message: "Avoid sourcing scripts from URLs. Download, verify integrity, then source from a local path to prevent supply-chain attacks.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1139")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1140(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid hash -r",
			input:    `hash -r`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid hash for existence check",
			input: `hash git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1140",
					Message: "Use `command -v cmd` instead of `hash cmd` for command existence checks. `command -v` provides clearer semantics in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1140")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1141(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid curl to file",
			input:    `curl -o file.tar.gz https://example.com/file.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid curl -sSL",
			input: `curl -sSL https://example.com/install.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1141",
					Message: "Avoid `curl -s URL | sh`. Download the script first, verify its integrity, then execute. Piping directly from the internet is a supply-chain risk.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1141")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1142(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid single grep",
			input:    `grep -E "foo|bar" file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep | grep",
			input: `grep foo file | grep bar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1142",
					Message: "Avoid chaining `grep | grep`. Combine into a single `grep -E` with alternation or use `awk` for multi-pattern matching to reduce pipeline processes.",
					Line:    1,
					Column:  15,
				},
			},
		},
		{
			name:     "non-pipe operator",
			input:    `echo hello && echo world`,
			expected: []katas.Violation{},
		},
		{
			name:     "pipe but not grep",
			input:    `cat file | sort`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1142")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1143(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid setopt",
			input:    `set -u`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid set -e",
			input: `set -e`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1143",
					Message: "Avoid `set -e`. It has surprising behavior with conditionals and subshells in Zsh. Use explicit error handling with `cmd || return 1` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1143")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1144(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid trap with name",
			input:    `trap cleanup EXIT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid trap with 0",
			input:    `trap cleanup 0`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid trap with number",
			input: `trap cleanup 15`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1144",
					Message: "Use signal names (`SIGTERM`, `SIGINT`, `EXIT`) instead of numbers in `trap`. Signal numbers vary across platforms.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1144")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1145(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid tr with complex set",
			input:    `tr -d '[:space:]'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid tr -d simple char",
			input: `tr -d ' '`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1145",
					Message: "Use `${var//char/}` instead of piping through `tr -d`. Parameter expansion is faster for simple character deletion.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1145")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1146(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid awk with file",
			input:    `awk '{print}' file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cat with flags piped",
			input:    `cat -n file | awk '{print}'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat | awk",
			input: `cat data.csv | awk -F, '{print $1}'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1146",
					Message: "Pass the file directly to `awk` instead of `cat file | awk`. Most text-processing tools accept file arguments.",
					Line:    1,
					Column:  14,
				},
			},
		},
		{
			name:  "invalid cat | sort",
			input: `cat names.txt | sort`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1146",
					Message: "Pass the file directly to `sort` instead of `cat file | sort`. Most text-processing tools accept file arguments.",
					Line:    1,
					Column:  15,
				},
			},
		},
		{
			name:     "non-pipe operator",
			input:    `echo hello && echo world`,
			expected: []katas.Violation{},
		},
		{
			name:     "pipe but left not cat",
			input:    `echo hello | sort`,
			expected: []katas.Violation{},
		},
		{
			name:     "cat no args piped",
			input:    `cat | awk '{print}'`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1146")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1147(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid mkdir -p",
			input:    `mkdir -p /tmp/a/b/c`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid mkdir single level",
			input:    `mkdir newdir`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid mkdir nested without -p",
			input: `mkdir /tmp/a/b/c`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1147",
					Message: "Use `mkdir -p` when creating nested directories. Without `-p`, `mkdir` fails if parent directories don't exist.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1147")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1148(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid compdef",
			input:    `compdef _git git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid compctl",
			input: `compctl -K _my_func mycommand`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1148",
					Message: "Use `compdef` instead of `compctl`. The `compctl` system is deprecated; use `compinit` and `compdef` for modern Zsh completions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1148")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1149(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid echo normal message",
			input:    `echo "Processing..."`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid echo error to stdout",
			input: `echo "Error: file not found"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1149",
					Message: "Error messages should go to stderr. Use `print -u2` or append `>&2` to separate error output from normal stdout.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1149")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1151(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid cat with file",
			input:    `cat file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat -A",
			input: `cat -A file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1151",
					Message: "Avoid `cat -A` for inspecting non-printable characters. Use `od -c` or `hexdump -C` for reliable cross-platform output.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1151")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1152(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -E",
			input:    `grep -E "pattern" file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep -P",
			input: `grep -P "\d+" file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1152",
					Message: "Avoid `grep -P` — it's unavailable on macOS. Use `zmodload zsh/pcre` with `pcre_compile`/`pcre_match` or `grep -E` for portable regex matching.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1152")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1153(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid diff for viewing",
			input:    `diff file1 file2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid diff -q for equality",
			input: `diff -q file1 file2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1153",
					Message: "Use `cmp -s file1 file2` instead of `diff -q`. `cmp -s` is faster for equality checks as it stops at the first difference.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1153")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1154(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid find without exec",
			input:    `find . -name "*.go"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid find -print0",
			input:    `find . -print0`,
			expected: []katas.Violation{},
		},
		{
			name:     "not find command",
			input:    `ls -la`,
			expected: []katas.Violation{},
		},
		{
			name:     "find with -exec ending in +",
			input:    `find . -exec rm {} +`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1154")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1155(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid whence -a",
			input:    `whence -a git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid which -a",
			input: `which -a git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1155",
					Message: "Use `whence -a` instead of `which -a`. Zsh `whence` is a reliable builtin for listing all command locations.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1155")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1156(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ln -s",
			input:    `ln -s target link`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid ln -sf",
			input:    `ln -sf target link`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid hard link",
			input: `ln target link`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1156",
					Message: "Use `ln -s` for symbolic links instead of hard links. Hard links share inodes and don't work across filesystems.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1156")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1157(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid strings with flags",
			input:    `strings -a binary`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid simple strings",
			input: `strings binary`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1157",
					Message: "Consider Zsh parameter expansion for string extraction from variables. `strings` is typically needed only for binary file analysis.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1157")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1158(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid chown without -R",
			input:    `chown user:group file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid chown -Rh",
			input:    `chown -Rh user:group dir`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid chown -R without -h",
			input: `chown -R user:group dir`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1158",
					Message: "Use `chown -Rh` or `chown -R --no-dereference` to prevent following symlinks during recursive ownership changes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1158")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1159(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid tar czf",
			input:    `tar czf archive.tar.gz dir`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid tar extract",
			input:    `tar xf archive.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid tar cf without compression",
			input: `tar cf archive.tar dir`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1159",
					Message: "Specify an explicit compression flag (`-z`, `-j`, `-J`) when creating tar archives. Relying on auto-detection reduces clarity and portability.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1159")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1160(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid curl",
			input:    `curl -o file.tar.gz https://example.com/file.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid wget",
			input: `wget https://example.com/file.tar.gz`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1160",
					Message: "Prefer `curl` over `wget` for portability. `curl` is pre-installed on macOS and most Linux distributions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1160")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1161(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sha256sum with file",
			input:    `sha256sum file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sha256sum in pipeline",
			input: `sha256sum -`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1161",
					Message: "Consider `zmodload zsh/sha256` or `zmodload zsh/md5` for hash operations. Zsh modules avoid spawning external hashing processes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1161")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1162(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid cp -a",
			input:    `cp -a src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cp single file",
			input:    `cp file.txt backup.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cp -r",
			input: `cp -r src/ dst/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1162",
					Message: "Use `cp -a` instead of `cp -r` to preserve permissions, timestamps, and symlinks. Archive mode ensures a faithful copy.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1162")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1163(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -m 1",
			input:    `grep -m 1 pattern file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep | head -1",
			input: `grep pattern file | head -1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1163",
					Message: "Use `grep -m 1` instead of `grep | head -1`. The `-m` flag stops after the first match without a pipeline.",
					Line:    1,
					Column:  19,
				},
			},
		},
		{
			name:     "non-pipe operator",
			input:    `echo hello && echo world`,
			expected: []katas.Violation{},
		},
		{
			name:     "pipe but left not grep",
			input:    `cat file | head -1`,
			expected: []katas.Violation{},
		},
		{
			name:     "grep piped to non-head",
			input:    `grep pattern file | sort`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1163")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1164(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sed substitution",
			input:    `sed -n 's/foo/bar/p'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sed with file",
			input:    `sed -n '3p' file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sed -n Np in pipeline",
			input: `sed -n '5p'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1164",
					Message: "Use Zsh array subscript `${lines[N]}` instead of `sed -n 'Np'`. Split input with `${(f)...}` then index directly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1164")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1165(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid awk with complex script",
			input:    `awk '{sum+=$1} END{print sum}'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid awk with file",
			input:    `awk '{print $1}' file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid simple awk print $1",
			input: `awk '{print $1}'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1165",
					Message: "Use Zsh parameter expansion (`${var%% *}` or `${var##* }`) instead of `awk '{print $1}'` for simple field extraction without spawning awk.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1165")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1166(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -i with file",
			input:    `grep -i pattern file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep -i in pipeline",
			input: `grep -i pattern`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1166",
					Message: "Use Zsh `(#i)` glob flag for case-insensitive matching instead of piping through `grep -i`. Example: `[[ $var == (#i)pattern ]]`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1166")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1167(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid timeout command",
			input: `timeout 5 cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1167",
					Message: "Avoid `timeout` — it's unavailable on macOS. Use Zsh `TMOUT` variable or `zmodload zsh/sched` for portable timeout functionality.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1167")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1168(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid readarray",
			input: `readarray -t arr`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1168",
					Message: "Use Zsh `${(f)$(cmd)}` instead of `readarray`. `readarray`/`mapfile` are Bash builtins not available in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid mapfile",
			input: `mapfile -t lines`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1168",
					Message: "Use Zsh `${(f)$(cmd)}` instead of `mapfile`. `readarray`/`mapfile` are Bash builtins not available in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1168")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1169(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid install without -m",
			input:    `install -d /usr/local/bin`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid install -m",
			input: `install -m 755 script /usr/local/bin/script`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1169",
					Message: "Consider using `cp` + `chmod` instead of `install -m`. Separate commands are clearer in shell scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1169")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1170(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid pushd -q",
			input:    `pushd -q /tmp`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid pushd without -q",
			input: `pushd /tmp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1170",
					Message: "Use `pushd -q` to suppress directory stack output in scripts. Without `-q`, the stack is printed on every call.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid popd without -q",
			input: `popd /tmp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1170",
					Message: "Use `popd -q` to suppress directory stack output in scripts. Without `-q`, the stack is printed on every call.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1170")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1171(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid print",
			input:    `print "hello\nworld"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid echo without -e",
			input:    `echo "hello"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid echo -e",
			input: `echo -e "hello\nworld"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1171",
					Message: "Use `print` instead of `echo -e`. Zsh `print` natively interprets escape sequences and is more portable than `echo -e`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1171")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1172(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid read -A",
			input:    `read -A arr`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid read -r",
			input:    `read -r line`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid read -a (Bash syntax)",
			input: `read -a arr`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1172",
					Message: "Use `read -A` instead of `read -a` in Zsh. The `-a` flag is Bash syntax; Zsh uses `-A` to read into arrays.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1172")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1173(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid column -s",
			input:    `column -s: file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid column -t",
			input: `column -t file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1173",
					Message: "Use Zsh `print -C N` for columnar output instead of `column -t`. The `print` builtin formats columns without spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1173")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1174(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid paste with files",
			input:    `paste -sd, file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid paste -sd in pipeline",
			input: `paste -s -d,`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1174",
					Message: "Use Zsh `${(j:delim:)array}` to join array elements instead of `paste -sd`. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1174")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1175(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid tput cols",
			input:    `tput cols`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid tput setaf",
			input: `tput setaf 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1175",
					Message: "Use Zsh `%F{color}` / `%f` or `$fg[color]` / `$reset_color` instead of `tput`. Zsh handles ANSI colors natively without spawning external processes.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid tput sgr0",
			input: `tput sgr0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1175",
					Message: "Use Zsh `%F{color}` / `%f` or `$fg[color]` / `$reset_color` instead of `tput`. Zsh handles ANSI colors natively without spawning external processes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1175")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1176(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid zparseopts",
			input:    `zparseopts -D -E -- v=verbose h=help`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid getopt",
			input: `getopt -o vh -l verbose,help`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1176",
					Message: "Use Zsh `zparseopts` instead of `getopt`. `zparseopts` supports long options, arrays, and is the native Zsh approach.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid getopts",
			input: `getopts "vh" opt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1176",
					Message: "Use Zsh `zparseopts` instead of `getopts`. `zparseopts` supports long options, arrays, and is the native Zsh approach.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1176")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1177(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid id -g",
			input:    `id -g`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid id -u",
			input: `id -u`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1177",
					Message: "Use `$UID` or `$EUID` instead of `id -u`. Zsh provides these as built-in variables.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1177")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1178(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid stty raw",
			input:    `stty raw`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid stty size",
			input: `stty size`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1178",
					Message: "Use `$COLUMNS` and `$LINES` instead of `stty size`. Zsh tracks terminal dimensions as built-in variables.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1178")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1179(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid date without format",
			input:    `date -u`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid date with format",
			input: `date '+%Y-%m-%d'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1179",
					Message: "Use `strftime` (via `zmodload zsh/datetime`) instead of `date +%Y-%m-%d`. Zsh date formatting avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1179")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1180(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid pgrep with -u flag",
			input:    `pgrep -u root sshd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid simple pgrep",
			input: `pgrep myprocess`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1180",
					Message: "For own background jobs, use Zsh job control (`jobs`, `kill %N`) instead of `pgrep`. Job control is more precise for script-spawned processes.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid pkill",
			input: `pkill -f myserver`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1180",
					Message: "For own background jobs, use Zsh job control (`jobs`, `kill %N`) instead of `pkill`. Job control is more precise for script-spawned processes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1180")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1181(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid open with -a flag",
			input:    `open -a Safari`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid xdg-open",
			input: `xdg-open https://example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1181",
					Message: "Use `$BROWSER` or check `$OSTYPE` instead of `xdg-open` for portable URL/file opening across Linux and macOS.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid open URL",
			input: `open https://example.com`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1181",
					Message: "Use `$BROWSER` or check `$OSTYPE` instead of `open` for portable URL/file opening across Linux and macOS.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1181")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1182(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid nc usage",
			input: `nc localhost 8080`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1182",
					Message: "Avoid `nc` for network operations in scripts. Use `curl` for HTTP or `zmodload zsh/net/tcp` for raw TCP connections with TLS support.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid netcat usage",
			input: `netcat -l 9090`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1182",
					Message: "Avoid `netcat` for network operations in scripts. Use `curl` for HTTP or `zmodload zsh/net/tcp` for raw TCP connections with TLS support.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1182")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1183(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ls -la",
			input:    `ls -la`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ls -t",
			input: `ls -t`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1183",
					Message: "Use Zsh glob qualifiers `*(om[1])` for newest file or `*(Om[1])` for oldest instead of `ls -t`. Glob qualifiers avoid spawning external processes.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid ls -ltr",
			input: `ls -ltr`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1183",
					Message: "Use Zsh glob qualifiers `*(om[1])` for newest file or `*(Om[1])` for oldest instead of `ls -t`. Glob qualifiers avoid spawning external processes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1183")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1184(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid diff without -u",
			input:    `diff file1 file2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid diff -u",
			input: `diff -u old.txt new.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1184",
					Message: "Consider `git diff` instead of `diff -u` when working in a repository. `git diff` provides better context and integration.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1184")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1185(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid wc -w with file",
			input:    `wc -w file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid wc -l",
			input:    `wc -l`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid wc -w in pipeline",
			input: `wc -w`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1185",
					Message: "Use Zsh `${#${(z)var}}` for word counting instead of piping through `wc -w`. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1185")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1186(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid unset -v",
			input:    `unset -v myvar`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid unset -f",
			input:    `unset -f myfunc`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid bare unset",
			input: `unset myvar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1186",
					Message: "Use `unset -v name` for variables or `unset -f name` for functions. Bare `unset` is ambiguous about what is being removed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1186")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1187(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid notify-send",
			input: `notify-send "Build complete"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1187",
					Message: "Wrap `notify-send` with an `$OSTYPE` check or `command -v` guard. It is Linux-only and will fail silently on macOS.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1187")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1188(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid export other var",
			input:    `export EDITOR=vim`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid export PATH",
			input: `export PATH=$PATH:/usr/local/bin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1188",
					Message: "Use `path+=(dir)` instead of `export PATH=$PATH:dir`. Zsh ties the `path` array to `$PATH` for cleaner manipulation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1188")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1189(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid source file",
			input:    `source /etc/profile`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid source /dev/stdin",
			input: `source /dev/stdin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1189",
					Message: "Avoid `source /dev/stdin`. Use `eval \"$(cmd)\"` for direct evaluation. `/dev/stdin` sourcing is fragile across platforms.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1189")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1190(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -v -e combined",
			input:    `grep -v -e foo -e bar file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid chained grep -v",
			input: `grep -v foo file | grep -v bar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1190",
					Message: "Combine `grep -v p1 | grep -v p2` into `grep -v -e p1 -e p2`. A single invocation avoids an unnecessary pipeline.",
					Line:    1,
					Column:  18,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1190")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1191(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid other command",
			input:    `reset`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid bare clear",
			input: `clear`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1191",
					Message: "Use `print -n '\\e[2J\\e[H'` instead of `clear`. ANSI escape sequences avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1191")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1192(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sleep 1",
			input:    `sleep 1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sleep 0",
			input: `sleep 0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1192",
					Message: "Remove `sleep 0` — it spawns a process that does nothing. Use `:` if an explicit no-op is needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1192")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1193(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid rm -f",
			input:    `rm -f file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid rm -i",
			input: `rm -i file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1193",
					Message: "Avoid `rm -i` in scripts — it prompts interactively and will hang in non-interactive execution. Remove `-i` or use explicit checks instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1193")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1194(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid single sed expression",
			input:    `sed -e 's/foo/bar/' file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid multiple sed -e",
			input: `sed -e 's/foo/bar/' -e 's/baz/qux/' file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1194",
					Message: "Combine multiple `sed -e` expressions into a single script: `sed 's/a/b/; s/c/d/'` is cleaner than multiple `-e` flags.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1194")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1195(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid umask 022",
			input:    `umask 022`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid umask 000",
			input: `umask 000`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1195",
					Message: "Avoid `umask 000` — it creates world-writable files. Use `umask 022` or `umask 077` for secure default permissions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1195")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1196(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid less with flag",
			input:    `less -R file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid less in script",
			input: `less file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1196",
					Message: "Avoid `less` in scripts — it requires interactive terminal input. Use `cat` or redirect output to a pager only when `$TERM` is available.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1196")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1197(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid more with flag",
			input:    `more -d file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid more in script",
			input: `more output.log`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1197",
					Message: "Avoid `more` in scripts — it requires an interactive terminal. Use `cat` for output or check `[[ -t 1 ]]` before paging.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1197")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1198(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sed -i",
			input:    `sed -i 's/old/new/' file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid nano in script",
			input: `nano config.yaml`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1198",
					Message: "Avoid `nano` in scripts — interactive editors hang without a terminal. Use `sed -i` or `ed` for scripted file editing.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid vim in script",
			input: `vim /etc/hosts`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1198",
					Message: "Avoid `vim` in scripts — interactive editors hang without a terminal. Use `sed -i` or `ed` for scripted file editing.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1198")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1199(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid curl check",
			input:    `curl -s http://localhost:8080`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid telnet",
			input: `telnet localhost 8080`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1199",
					Message: "Avoid `telnet` in scripts — it is interactive and insecure. Use `curl` for HTTP checks or `zmodload zsh/net/tcp` for port testing.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1199")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
