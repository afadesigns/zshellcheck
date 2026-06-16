// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1001(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid array access",
			input:    `echo ${my_array[1]}`,
			expected: []katas.Violation{},
		},
		{
			name:  "braced form preferred for array access",
			input: `echo $my_array[1]`,
			expected: []katas.Violation{
				{
					KataID: "ZC1001",
					Message: "Prefer `${...}` for array element access. " +
						"`$my_array[...]` is valid Zsh, but the braced form is unambiguous and robust under `KSH_ARRAYS`.",
					Line:   1,
					Column: 6,
				},
			},
		},
		{
			// `$+commands[ls]` is the parameter-existence operator, not
			// array element access — it returns 0 or 1 inside `(( … ))`.
			name:     "existence test on commands is not element access",
			input:    `(( $+commands[ls] ))`,
			expected: []katas.Violation{},
		},
		{
			name:     "existence test on functions is not element access",
			input:    `(( $+functions[foo] ))`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1001")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1002(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid command substitution",
			input:    `x=$(ls)`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid backticks",
			input: `x=` + "`ls`",
			expected: []katas.Violation{
				{
					KataID: "ZC1002",
					Message: "Use $(...) instead of backticks for command substitution. " +
						"The `$(...)` syntax is more readable and can be nested easily.",
					Line:   1,
					Column: 3,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1002")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1003(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid arithmetic test",
			input:    "(( 1 > 0 ))",
			expected: []katas.Violation{},
		},
		{
			name:  "invalid arithmetic test",
			input: "[ 1 -gt 0 ]", // Ensure spaces
			expected: []katas.Violation{
				{
					KataID:  "ZC1003",
					Message: "Use `((...))` for arithmetic comparisons instead of `[` or `test`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1003")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1004(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid return",
			input:    `my_func() { return 0; }`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid exit in subshell",
			input:    `my_func() { ( exit 1 ) }`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid exit in command sub",
			input:    `my_func() { local x=$(exit 1); }`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid exit",
			input: `my_func() { exit 1; }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1004",
					Message: "Use `return` instead of `exit` in functions to avoid killing the shell.",
					Line:    1,
					Column:  13,
				},
			},
		},
		{
			name:  "invalid exit deep",
			input: `my_func() { if true; then exit 1; fi }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1004",
					Message: "Use `return` instead of `exit` in functions to avoid killing the shell.",
					Line:    1,
					Column:  27,
				},
			},
		},
		{
			name:  "exit in function keyword style",
			input: `function my_func { exit 1; }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1004",
					Message: "Use `return` instead of `exit` in functions to avoid killing the shell.",
					Line:    1,
					Column:  20,
				},
			},
		},
		{
			name:     "non-function node",
			input:    `exit 0`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1004")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1005(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `which ls`,
			expected: []katas.Violation{
				{
					KataID: "ZC1005",
					Message: "Use whence instead of which. The `whence` command is a built-in Zsh command " +
						"that provides a more reliable and consistent way to find the location of a command.",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			input:    `whence ls`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1005")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1006(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `test 1 -eq 1`,
			expected: []katas.Violation{
				{
					KataID: "ZC1006",
					Message: "Prefer [[ over test for tests. " +
						"[[ is a Zsh keyword that offers safer and more powerful conditional expressions.",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			input:    `[[ 1 -eq 1 ]]`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1006")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1007(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `chmod 777 file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1007",
					Message: "Avoid using `chmod 777`. It is a security risk.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `chmod 755 file.txt`,
			expected: []katas.Violation{},
		},
		{
			input:    `ls -l`,
			expected: []katas.Violation{},
		},
		{
			input: `chmod 777 file1 file2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1007",
					Message: "Avoid using `chmod 777`. It is a security risk.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1007")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1008(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `let x = 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1008",
					Message: "Use `\\$(())` for arithmetic operations instead of `let`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `x=1`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1008")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1009(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `expr 1 + 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1009",
					Message: "Use `((...))` for C-style arithmetic instead of `expr`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `echo "hello"`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1009")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1010(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			// An arithmetic comparison defers to ZC1003 (`(( … ))`), not
			// ZC1010 (`[[ … ]]`), so ZC1010 does not fire here.
			input:    `[ 1 -eq 1 ]`,
			expected: []katas.Violation{},
		},
		{
			// A non-arithmetic test still gets the `[[ … ]]` advice.
			input: `[ -f file ]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1010",
					Message: "Use `[[ ... ]]` instead of `[ ... ]` or `test`. `[[` is safer and more powerful.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `[[ 1 -eq 1 ]]`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1010")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1011(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `git rev-parse HEAD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1011",
					Message: "Avoid using `git` plumbing commands in scripts. They are not guaranteed to be stable.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `git branch`,
			expected: []katas.Violation{},
		},
		{
			input:    `ls -l`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1011")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1012(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "read without flags",
			input: `read line`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1012",
					Message: "Use `read -r` to read input without interpreting backslashes.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "read with -r",
			input:    `read -r line`,
			expected: []katas.Violation{},
		},
		{
			name:     "read with -er",
			input:    `read -er line`, // heuristic support
			expected: []katas.Violation{},
		},
		{
			name:     "read -k reads raw chars so -r is a no-op",
			input:    `read -k 3 key`,
			expected: []katas.Violation{},
		},
		{
			name:     "read -q reads a single y/n char",
			input:    `read -q "REPLY?ok"`,
			expected: []katas.Violation{},
		},
		{
			name:     "read -z reads the editor buffer",
			input:    `read -z buf`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1012")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1013(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `let x = 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1013",
					Message: "Use `((...))` for arithmetic operations instead of `let`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `[[ -f file ]]`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1013")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1014(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `git checkout my-branch`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1014",
					Message: "Use `git switch` or `git restore` instead of the ambiguous `git checkout`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `[ -f file ]`,
			expected: []katas.Violation{},
		},
		{
			input:    `git restore my-file`,
			expected: []katas.Violation{},
		},
		{
			input:    `git commit -m "message"`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1014")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1015(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: "`ls`",
			expected: []katas.Violation{
				{
					KataID:  "ZC1015",
					Message: "Use `$(...)` for command substitution instead of backticks.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `$(ls)`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1015")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1016(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "safe read",
			input:    `read name`,
			expected: []katas.Violation{},
		},
		{
			name:     "safe read password with -s",
			input:    `read -s password`,
			expected: []katas.Violation{},
		},
		{
			name:     "safe read with combined flags",
			input:    `read -rs password`,
			expected: []katas.Violation{},
		},
		{
			name:  "unsafe read password",
			input: `read password`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1016",
					Message: "Use `read -s` to hide input when reading sensitive variable 'password'.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "unsafe read with prompt",
			input: `read "secret_key?Enter key: "`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1016",
					Message: "Use `read -s` to hide input when reading sensitive variable 'secret_key'.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "unsafe read multiple vars",
			input: `read user password`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1016",
					Message: "Use `read -s` to hide input when reading sensitive variable 'password'.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			// `key` is a substring of `monkey`; whole-name matching
			// must not flag it.
			name:     "non-secret substring monkey",
			input:    `read monkey`,
			expected: []katas.Violation{},
		},
		{
			name:     "non-secret substring keyword",
			input:    `read keyword`,
			expected: []katas.Violation{},
		},
		{
			name:     "non-secret substring keys",
			input:    `read keys`,
			expected: []katas.Violation{},
		},
		{
			// `pwd` was an over-broad substring; the print-working-dir
			// name is not a secret.
			name:     "non-secret pwd",
			input:    `read pwd`,
			expected: []katas.Violation{},
		},
		{
			// `-s` is meaningless when stdin comes from a file: no
			// terminal echo to hide.
			name:     "secret from file redirection",
			input:    `read secret < /etc/secret`,
			expected: []katas.Violation{},
		},
		{
			name:     "secret from here-string",
			input:    `read password <<< "$data"`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1016")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1017(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			// No backslash escape — nothing for `print` to interpret,
			// so `-r` makes no difference and the kata stays quiet.
			input:    `print "hello"`,
			expected: []katas.Violation{},
		},
		{
			input: `print "a\tb"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1017",
					Message: "Use `print -r` to print strings literally.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `print -r "a\tb"`,
			expected: []katas.Violation{},
		},
		{
			// Prompt escapes are expanded by `-P`; `-r` does not suppress
			// them and the content carries no backslash escape.
			input:    `print -P "%F{green}done%f"`,
			expected: []katas.Violation{},
		},
		{
			input:    `print -l -- $usage`,
			expected: []katas.Violation{},
		},
		{
			input:    `print "${cmd}"`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1017")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

// ZC1018 was retired as a duplicate of ZC1009. Kept as a no-op stub so
// legacy `disabled_katas` lists keep parsing; the detection runs under
// ZC1009 now.

func TestZC1018Stub(t *testing.T) {
	cases := []string{
		"echo hi",
		"ls",
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			v := testutil.Check(in, "ZC1018")
			testutil.AssertViolations(t, in, v, []katas.Violation{})
		})
	}
}

// ZC1019 was retired as a duplicate of ZC1005. Kept as a no-op stub so
// legacy `disabled_katas` lists keep parsing; the detection runs under
// ZC1005 now.

func TestZC1019Stub(t *testing.T) {
	cases := []string{
		"echo hi",
		"ls",
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			v := testutil.Check(in, "ZC1019")
			testutil.AssertViolations(t, in, v, []katas.Violation{})
		})
	}
}

func TestCheckZC1020(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `test 1 -eq 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1020",
					Message: "Use `[[ ... ]]` for tests instead of `test`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `[[ 1 -eq 1 ]]`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1020")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1021(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `chmod 755 file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1021",
					Message: "Use symbolic permissions with `chmod` instead of octal.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `chmod u+x file`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1021")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1022(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `let x=1+1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1022",
					Message: "Use `$((...))` for arithmetic expansion instead of `let`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `x=$((1+1))`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1022")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

// ZC1023 was retired as a duplicate of ZC1022 (see issue #345). It is
// kept as a no-op stub so legacy `disabled_katas` lists that reference
// it keep parsing; the canonical `let` → `$((...))` guidance fires
// under ZC1022 now.

func TestCheckZC1023(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{input: "let x=1+1", expected: []katas.Violation{}},
		{input: "x=$((1+1))", expected: []katas.Violation{}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1023")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

// ZC1024 was retired as a duplicate of ZC1022 (see issue #345). It is
// kept as a no-op stub so legacy `disabled_katas` lists that reference
// it keep parsing; the canonical `let` → `$((...))` guidance fires
// under ZC1022 now.

func TestCheckZC1024(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{input: "let x=1+1", expected: []katas.Violation{}},
		{input: "x=$((1+1))", expected: []katas.Violation{}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1024")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

// ZC1025 was retired as a duplicate of ZC1022 (see issue #345). It is
// kept as a no-op stub so legacy `disabled_katas` lists that reference
// it keep parsing; the canonical `let` → `$((...))` guidance fires
// under ZC1022 now.

func TestCheckZC1025(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{input: "let x=1+1", expected: []katas.Violation{}},
		{input: "x=$((1+1))", expected: []katas.Violation{}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1025")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

// ZC1026 was retired as a duplicate of ZC1022 (see issue #345). It is
// kept as a no-op stub so legacy `disabled_katas` lists that reference
// it keep parsing; the canonical `let` → `$((...))` guidance fires
// under ZC1022 now.

func TestCheckZC1026(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{input: "let x=1+1", expected: []katas.Violation{}},
		{input: "x=$((1+1))", expected: []katas.Violation{}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1026")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

// ZC1027 was retired as a duplicate of ZC1022 (see issue #345). It is
// kept as a no-op stub so legacy `disabled_katas` lists that reference
// it keep parsing; the canonical `let` → `$((...))` guidance fires
// under ZC1022 now.

func TestCheckZC1027(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{input: "let x=1+1", expected: []katas.Violation{}},
		{input: "x=$((1+1))", expected: []katas.Violation{}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1027")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

// ZC1028 was retired as a duplicate of ZC1022 (see issue #345). It is
// kept as a no-op stub so legacy `disabled_katas` lists that reference
// it keep parsing; the canonical `let` → `$((...))` guidance fires
// under ZC1022 now.

func TestCheckZC1028(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{input: "let x=1+1", expected: []katas.Violation{}},
		{input: "x=$((1+1))", expected: []katas.Violation{}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1028")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

// ZC1029 was retired as a duplicate of ZC1022 (see issue #345). It is
// kept as a no-op stub so legacy `disabled_katas` lists that reference
// it keep parsing; the canonical `let` → `$((...))` guidance fires
// under ZC1022 now.

func TestCheckZC1029(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{input: "let x=1+1", expected: []katas.Violation{}},
		{input: "x=$((1+1))", expected: []katas.Violation{}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1029")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1030(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			// ZC1030's `printf` advice is consolidated under ZC1037
			// (`print -r --`), so it no longer fires on `echo`.
			name:     "echo with a simple string defers to ZC1037",
			input:    `echo "hello"`,
			expected: []katas.Violation{},
		},
		{
			name:     "printf with a simple string",
			input:    `printf "hello"`,
			expected: []katas.Violation{},
		},
		{
			name:     "echo with a variable defers to ZC1037",
			input:    `echo "$foo"`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1030")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1031(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `#!/bin/zsh
echo "hello"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1031",
					Message: "Use `#!/usr/bin/env zsh` for portability instead of `#!/bin/zsh`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input: `#!/usr/bin/env zsh
echo "hello"`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1031")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1032(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `let i=i+1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1032",
					Message: "Use `(( i++ ))` for C-style incrementing instead of `let i=i+1`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `(( i++ ))`,
			expected: []katas.Violation{},
		},
		{
			input:    `let i=j+1`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1032")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

// ZC1033 was retired as a duplicate of ZC1022 (see issue #345). It is
// kept as a no-op stub so legacy `disabled_katas` lists that reference
// it keep parsing; the canonical `let` → `$((...))` guidance fires
// under ZC1022 now.

func TestCheckZC1033(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{input: "let x=1+1", expected: []katas.Violation{}},
		{input: "x=$((1+1))", expected: []katas.Violation{}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1033")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1034(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `which ls`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1034",
					Message: "Use `command -v` instead of `which` for portability.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `command -v ls`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1034")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

// ZC1035 was retired as a duplicate of ZC1022 (see issue #345). It is
// kept as a no-op stub so legacy `disabled_katas` lists that reference
// it keep parsing; the canonical `let` → `$((...))` guidance fires
// under ZC1022 now.

func TestCheckZC1035(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{input: "let x=1+1", expected: []katas.Violation{}},
		{input: "x=$((1+1))", expected: []katas.Violation{}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1035")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1036(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `test -f file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1036",
					Message: "Prefer `[[ ... ]]` over `test` command for conditional expressions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1036")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1037(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid print",
			input:    `print -r -- "$var"`,
			expected: []katas.Violation{},
		},
		{
			// ZC1037 is now the single canonical echo recommendation and
			// fires on every `echo`, including literal ones.
			name:  "literal echo is flagged",
			input: `echo 'cost is $5'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1037",
					Message: "Use `print -r --` instead of `echo`. `echo` behaviour varies with options and escape handling; `print -r --` is the reliable Zsh builtin.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid echo with variable",
			input: `echo $var`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1037",
					Message: "Use `print -r --` instead of `echo`. `echo` behaviour varies with options and escape handling; `print -r --` is the reliable Zsh builtin.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid echo with quoted variable",
			input: `echo "$var"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1037",
					Message: "Use `print -r --` instead of `echo`. `echo` behaviour varies with options and escape handling; `print -r --` is the reliable Zsh builtin.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1037")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestCheckZC1038(t *testing.T) {
	tests := []struct {
		input    string
		expected []katas.Violation
	}{
		{
			input: `cat file | grep "foo"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1038",
					Message: "Avoid useless use of cat. Prefer `command file` or `command < file` over `cat file | command`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			input:    `cat file1 file2 | grep "foo"`, // Valid concatenation
			expected: []katas.Violation{},
		},
		{
			input:    `grep "foo" file`, // Direct file access
			expected: []katas.Violation{},
		},
		{
			input:    `cat | grep "foo"`, // Reading from stdin
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1038")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1039(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid usage",
			input:    `rm /tmp/file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid usage recursive",
			input:    `rm -rf /tmp/dir`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid rm root",
			input: `rm /`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1039",
					Message: "Avoid `rm` on the root directory `/`. This is highly dangerous.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:  "invalid rm root quoted",
			input: `rm "/"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1039",
					Message: "Avoid `rm` on the root directory `/`. This is highly dangerous.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:  "invalid rm root single quoted",
			input: `rm '/'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1039",
					Message: "Avoid `rm` on the root directory `/`. This is highly dangerous.",
					Line:    1,
					Column:  4,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1039")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1040(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "safe glob with nullglob qualifier",
			input:    `for f in *.txt(N); do echo $f; done`,
			expected: []katas.Violation{},
		},
		{
			// Nullglob `N` may appear with other qualifiers in any
			// order: `(@N)`, `(.N)` — all are nullglob-enabled.
			name:     "safe glob with combined N qualifier",
			input:    `for f in $x/*(@N); do echo $f; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "safe glob with dot-N qualifier",
			input:    `for f in *.log(.N); do echo $f; done`,
			expected: []katas.Violation{},
		},
		{
			// An array subscript `$arr[@]` is not a filename glob and
			// never triggers nomatch.
			name:     "array subscript at-all is not a glob",
			input:    `for w in $empty[@]; do echo $w; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "array subscript index is not a glob",
			input:    `for f in $arr[1]; do echo $f; done`,
			expected: []katas.Violation{},
		},
		{
			// A brace range expands in the shell and never nomatches.
			name:     "brace range is not a glob",
			input:    `for i in {3..1}; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "no glob pattern",
			input:    `for f in a b c; do echo $f; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "arithmetic for loop",
			input:    `for (( i=0; i<10; i++ )); do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "quoted string is not a glob",
			input:    `for f in "*.txt"; do echo $f; done`,
			expected: []katas.Violation{},
		},
		{
			name:  "bare glob without qualifier fires",
			input: `for f in *.txt; do echo $f; done`,
			expected: []katas.Violation{
				{
					KataID: "ZC1040",
					Message: "Glob pattern '*.txt' may error if no files match. " +
						"Append '(N)' to enable nullglob behavior: '*.txt(N)'",
					Line:   1,
					Column: 10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1040")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1041(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "safe printf with string literal",
			input:    `printf '%s\n' "hello"`,
			expected: []katas.Violation{},
		},
		{
			name:     "not printf command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "printf with no args",
			input:    `printf`,
			expected: []katas.Violation{},
		},
		{
			name:  "printf with variable as format",
			input: `printf $var`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1041",
					Message: "Do not use variables in printf format string. Use `printf '..%s..' \"$var\"` instead.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:     "printf with safe static format",
			input:    `printf 'hello world'`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1041")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1042(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "safe loop with $@",
			input:    `for i in "$@"; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "loop over plain items",
			input:    `for i in a b c; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "arithmetic loop",
			input:    `for (( i=0; i<10; i++ )); do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "loop with $* (parsed as separate tokens)",
			input:    `for i in $*; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "for-each with string literal items",
			input:    `for i in "one" "two"; do echo $i; done`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1042")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1043(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "no function definition",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "global assignment in function",
			input: "myfunc() { x=1; }",
			expected: []katas.Violation{
				{KataID: "ZC1043", Message: "Variable 'x' is assigned without 'local'. It will be global. Use `local x=1`.", Line: 1, Column: 12},
			},
		},
		{
			name:     "local declaration in function",
			input:    "myfunc() { local x=1; }",
			expected: []katas.Violation{},
		},
		{
			// REPLY/reply/MATCH/... are Zsh's value-return parameters;
			// they are intentionally global and must not be flagged.
			name:     "REPLY return parameter not flagged",
			input:    "myfunc() { REPLY=42 }",
			expected: []katas.Violation{},
		},
		{
			name:     "reply array return parameter not flagged",
			input:    "myfunc() { reply=(a b c) }",
			expected: []katas.Violation{},
		},
		{
			// ZLE widgets write editor state; `local BUFFER` would sever
			// the editor binding. The buffer specials are global by design.
			name:     "ZLE buffer specials not flagged",
			input:    "my-widget() { LBUFFER=\"$LBUFFER foo\"; CURSOR=3; BUFFER=bar }",
			expected: []katas.Violation{},
		},
		{
			name:     "region_highlight ZLE array not flagged",
			input:    "my-widget() { region_highlight=(\"0 5 bold\") }",
			expected: []katas.Violation{},
		},
		{
			// Prompt-setup functions set the prompt globally; `local PROMPT`
			// discards it on return.
			name:     "prompt strings not flagged",
			input:    "prompt_foo_setup() { PROMPT='%n'; RPROMPT='%t'; prompt_opts=(cr subst) }",
			expected: []katas.Violation{},
		},
		{
			// Regression for #1229 — empty-RHS assignment must not
			// panic during message build. Hint still emitted with an
			// empty RHS rendered in the template.
			name:  "empty-RHS assignment does not panic",
			input: "myfunc() { empty= }",
			expected: []katas.Violation{
				{
					KataID: "ZC1043",
					Message: "Variable 'empty' is assigned without 'local'. It will be global. " +
						"Use `local empty=`.",
					Line:   1,
					Column: 12,
				},
			},
		},
		{
			// Regression for #1332 — an inline env-var prefix scopes the
			// assignment to the following command; it is not a global.
			name:     "env-prefix on a command is not flagged",
			input:    "myfunc() { DEBUG=true echo \"foo\"; }",
			expected: []katas.Violation{},
		},
		{
			name:     "multiple env-prefix assignments are not flagged",
			input:    "myfunc() { A=1 B=2 mycmd; }",
			expected: []katas.Violation{},
		},
		{
			// A `;` ends the assignment, so it is a standalone global
			// and must still be flagged.
			name:  "semicolon-separated assignment is still flagged",
			input: "myfunc() { DEBUG=true; echo foo; }",
			expected: []katas.Violation{
				{KataID: "ZC1043", Message: "Variable 'DEBUG' is assigned without 'local'. It will be global. Use `local DEBUG=true`.", Line: 1, Column: 12},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1043")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1044(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "cd with error handling",
			input:    `cd /tmp || exit 1`,
			expected: []katas.Violation{},
		},
		{
			name:  "unchecked cd",
			input: `cd /tmp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1044",
					Message: "Use `cd ... || return` (or `exit`) in case cd fails.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "no cd command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "cd in if condition",
			input:    `if cd /tmp; then echo ok; fi`,
			expected: []katas.Violation{},
		},
		{
			name:  "cd with && chain",
			input: `cd /tmp && echo ok`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1044",
					Message: "Use `cd ... || return` (or `exit`) in case cd fails.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "cd negated in condition",
			input:    `if ! cd /tmp; then echo fail; fi`,
			expected: []katas.Violation{},
		},
		{
			name:  "cd in while loop body",
			input: `while true; do cd /tmp; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1044",
					Message: "Use `cd ... || return` (or `exit`) in case cd fails.",
					Line:    1,
					Column:  16,
				},
			},
		},
		{
			name:  "cd in for loop body",
			input: `for d in /tmp /var; do cd $d; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1044",
					Message: "Use `cd ... || return` (or `exit`) in case cd fails.",
					Line:    1,
					Column:  24,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1044")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1045(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "safe local declaration",
			input:    `local var`,
			expected: []katas.Violation{},
		},
		{
			name:     "regular command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "local with simple value",
			input:    `local var=hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "local with command substitution",
			input: `local var=$(date)`,
			expected: []katas.Violation{
				{
					KataID: "ZC1045",
					Message: "Declare and assign separately to avoid masking return values. " +
						"`local var=$(cmd)` masks the exit code of `cmd`.",
					Line:   1,
					Column: 7,
				},
			},
		},
		{
			name:  "readonly with command substitution",
			input: `readonly var=$(whoami)`,
			expected: []katas.Violation{
				{
					KataID: "ZC1045",
					Message: "Declare and assign separately to avoid masking return values. " +
						"`readonly var=$(cmd)` masks the exit code of `cmd`.",
					Line:   1,
					Column: 10,
				},
			},
		},
		{
			name:  "declare with command substitution",
			input: `declare var=$(date)`,
			expected: []katas.Violation{
				{
					KataID: "ZC1045",
					Message: "Declare and assign separately to avoid masking return values. " +
						"`declare var=$(cmd)` masks the exit code of `cmd`.",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name:     "echo is not local or readonly",
			input:    `echo $(date)`,
			expected: []katas.Violation{},
		},
		{
			name:     "arithmetic expansion masks no command",
			input:    `local var=$(( 1 + 1 ))`,
			expected: []katas.Violation{},
		},
		{
			name:     "arithmetic with name operand",
			input:    `local elapsed=$(( EPOCHSECONDS - start ))`,
			expected: []katas.Violation{},
		},
		{
			name:     "arithmetic modulo masks no command",
			input:    `readonly slot=$(( RANDOM % 1000 ))`,
			expected: []katas.Violation{},
		},
		{
			name:     "integer-literal arithmetic masks no command",
			input:    `local n=$(( 5 ))`,
			expected: []katas.Violation{},
		},
		{
			// A pipeline inside the substitution is a command, so the
			// exit code is masked and the kata fires.
			name:  "pipeline command substitution masks return",
			input: `local out=$(a | b)`,
			expected: []katas.Violation{
				{
					KataID: "ZC1045",
					Message: "Declare and assign separately to avoid masking return values. " +
						"`local var=$(cmd)` masks the exit code of `cmd`.",
					Line:   1,
					Column: 7,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1045")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1046(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "direct eval",
			input: `eval "echo hello"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1046",
					Message: "Avoid `eval`. It allows execution of arbitrary code and is hard to debug.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "not eval",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "builtin eval",
			input: `builtin eval "echo hello"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1046",
					Message: "Avoid `eval`. It allows execution of arbitrary code and is hard to debug.",
					Line:    1,
					Column:  9,
				},
			},
		},
		{
			name:  "command eval",
			input: `command eval "echo hello"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1046",
					Message: "Avoid `eval`. It allows execution of arbitrary code and is hard to debug.",
					Line:    1,
					Column:  9,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1046")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1047(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "sudo command",
			input: `sudo apt install vim`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1047",
					Message: "Avoid `sudo` in scripts. Run the entire script as root if privileges are required.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "not sudo",
			input:    `apt install vim`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1047")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1048(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "source with absolute path",
			input:    `source /etc/profile`,
			expected: []katas.Violation{},
		},
		{
			name:     "not source command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "source with no arguments",
			input:    `source`,
			expected: []katas.Violation{},
		},
		{
			name:  "source with relative path",
			input: `source ./lib.zsh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1048",
					Message: "Avoid `source` with relative paths. Use `${0:a:h}/...` to resolve relative to the script.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:  "dot with relative path",
			input: `. ../lib.zsh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1048",
					Message: "Avoid `source` with relative paths. Use `${0:a:h}/...` to resolve relative to the script.",
					Line:    1,
					Column:  3,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1048")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1049(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "alias usage",
			input: `alias ll='ls -la'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1049",
					Message: "Prefer functions over aliases. Aliases are expanded at parse time and can behave unexpectedly in scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "function instead of alias",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			// Global aliases have no function equivalent.
			name:     "global alias is not flagged",
			input:    `alias -g GP='| grep'`,
			expected: []katas.Violation{},
		},
		{
			// Suffix aliases have no function equivalent.
			name:     "suffix alias is not flagged",
			input:    `alias -s md=less`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1049")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1050(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "safe glob loop",
			input:    `for f in *.txt; do echo $f; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "arithmetic loop",
			input:    `for (( i=0; i<5; i++ )); do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:  "loop over ls output",
			input: `for f in $(ls); do echo $f; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1050",
					Message: "Avoid iterating over `ls` output. Use globs (e.g. `*.txt`) to handle filenames with spaces correctly.",
					Line:    1,
					Column:  10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1050")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1051(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "rm with quoted variable",
			input:    `rm "$file"`,
			expected: []katas.Violation{},
		},
		{
			name:     "rm with literal path",
			input:    `rm /tmp/foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "not rm command",
			input:    `echo $var`,
			expected: []katas.Violation{},
		},
		{
			name:  "rm with unquoted variable",
			input: `rm $file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1051",
					Message: "Unquoted expansion in `rm`. An empty value changes the target — `rm -rf $dir/` becomes `rm -rf /`. Guard with `${dir:?}`.",
					Line:    1,
					Column:  4,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1051")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1052(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "sed without -i",
			input:    `sed 's/foo/bar/' file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "not sed command",
			input:    `grep foo bar`,
			expected: []katas.Violation{},
		},
		{
			name:     "sed with other flags",
			input:    `sed -e 's/foo/bar/' file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "sed with -i",
			input: `sed -i 's/foo/bar/' file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1052",
					Message: "`sed -i` is non-portable (GNU vs BSD differences). Use `perl -i` or a temporary file.",
					Line:    1,
					Column:  5,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1052")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1053(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "grep -q in condition",
			input:    `if grep -q pattern file; then echo found; fi`,
			expected: []katas.Violation{},
		},
		{
			name:     "no grep in condition",
			input:    `if true; then echo yes; fi`,
			expected: []katas.Violation{},
		},
		{
			name:  "grep without -q in condition",
			input: `if grep pattern file; then echo found; fi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:     "grep -rq combined flags",
			input:    `if grep -rq pattern dir; then echo found; fi`,
			expected: []katas.Violation{},
		},
		{
			name:     "grep --quiet in condition",
			input:    `if grep --quiet pattern file; then echo found; fi`,
			expected: []katas.Violation{},
		},
		{
			name:  "egrep without -q in condition",
			input: `if egrep pattern file; then echo found; fi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:  "grep piped in condition",
			input: `if echo test | grep pattern; then echo found; fi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    1,
					Column:  16,
				},
			},
		},
		{
			name:  "grep without -q in while condition",
			input: `while grep pattern file; do echo loop; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:     "grep -q in while condition",
			input:    `while grep -q pattern file; do echo loop; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "grep --silent in condition",
			input:    `if grep --silent pattern file; then echo found; fi`,
			expected: []katas.Violation{},
		},
		{
			name:  "fgrep without -q",
			input: `if fgrep pattern file; then echo found; fi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:  "zgrep without -q",
			input: `if zgrep pattern file; then echo found; fi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:     "grep left side of pipe is silenced",
			input:    `if grep pattern file | wc -l; then echo found; fi`,
			expected: []katas.Violation{},
		},
		{
			name:     "non-grep command in condition",
			input:    `if ls /tmp; then echo ok; fi`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1053")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1054(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "no range pattern",
			input:    `grep foo bar`,
			expected: []katas.Violation{},
		},
		{
			name:     "command with no args",
			input:    `ls`,
			expected: []katas.Violation{},
		},
		{
			name:     "digit range is locale-invariant",
			input:    `grep -o '[0-9]*'`,
			expected: []katas.Violation{},
		},
		{
			name:     "octal digit range is locale-invariant",
			input:    `grep -o '[0-7]'`,
			expected: []katas.Violation{},
		},
		{
			name:  "alphabetic range is collation-sensitive",
			input: `grep -o '[a-z]*'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1054",
					Message: "Ranges like `[a-z]` are locale-dependent. Use POSIX classes like `[[:lower:]]` or `[[:digit:]]`.",
					Line:    1,
					Column:  9,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1054")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1055(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "no comparison",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "compare with empty double-quoted string",
			input: `[[ $x == "" ]]`,
			expected: []katas.Violation{
				{KataID: "ZC1055", Message: "Use `[[ -z ... ]]` instead of comparing with empty string.", Line: 1, Column: 8},
			},
		},
		{
			name:  "compare not equal empty string",
			input: `[[ $x != "" ]]`,
			expected: []katas.Violation{
				{KataID: "ZC1055", Message: "Use `[[ -n ... ]]` instead of comparing with empty string.", Line: 1, Column: 8},
			},
		},
		{
			name:     "compare with non-empty string",
			input:    `[[ $x == "hello" ]]`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1055")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1056(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "normal command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "command substitution as statement",
			input:    `$(echo hello)`,
			expected: []katas.Violation{},
		},
		{
			name:     "arithmetic command",
			input:    `((x + 1))`,
			expected: []katas.Violation{},
		},
		{
			name:     "variable expansion command",
			input:    `$var`,
			expected: []katas.Violation{},
		},
		{
			name:     "ls command no violation",
			input:    `ls -la`,
			expected: []katas.Violation{},
		},
		{
			name:     "assignment statement",
			input:    `x=5`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1056")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1057(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "no ls assignment",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "safe array assignment",
			input:    `files=(*)`,
			expected: []katas.Violation{},
		},
		{
			name:     "echo not assignment",
			input:    `echo $(ls)`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1057")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1058(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "sudo without redirection",
			input:    `sudo apt install vim`,
			expected: []katas.Violation{},
		},
		{
			name:     "not sudo command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1058")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1059(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "rm with literal path",
			input:    `rm /tmp/foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "not rm command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "rm with no arguments",
			input:    `rm`,
			expected: []katas.Violation{},
		},
		{
			name:  "rm with ${var} argument",
			input: `rm ${dir}`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1059",
					Message: "Use `${var:?}` or ensure the variable is set before using it in `rm`.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:     "rm with flags and literal",
			input:    `rm -rf /tmp/build`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1059")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1060(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "no pipe",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "pipe without ps",
			input:    `cat file | grep pattern`,
			expected: []katas.Violation{},
		},
		{
			name:  "ps piped to grep",
			input: `ps aux | grep myprocess`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1060",
					Message: "`ps | grep pattern` matches the grep process itself. Use `grep [p]attern` to exclude the grep process.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:     "ps piped to non-grep command",
			input:    `ps aux | sort`,
			expected: []katas.Violation{},
		},
		{
			name:     "non-pipe operator",
			input:    `echo a && echo b`,
			expected: []katas.Violation{},
		},
		{
			name:  "ps piped to grep with flag args only",
			input: `ps aux | grep -i myprocess`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1060",
					Message: "`ps | grep pattern` matches the grep process itself. Use `grep [p]attern` to exclude the grep process.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:     "non-ps command piped to grep",
			input:    `ls -la | grep foo`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1060")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1061(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "seq usage",
			input: `seq 1 10`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1061",
					Message: "Prefer `{start..end}` range expansion over `seq`. It is built-in and faster.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "no seq",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1061")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1062(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "egrep usage",
			input: `egrep 'foo|bar' file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1062",
					Message: "`egrep` is deprecated. Use `grep -E` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "grep -E usage",
			input:    `grep -E 'foo|bar' file`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1062")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1063(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "fgrep usage",
			input: `fgrep 'literal' file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1063",
					Message: "`fgrep` is deprecated. Use `grep -F` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "grep -F usage",
			input:    `grep -F 'literal' file`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1063")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1064(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "type command",
			input: `type ls`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1064",
					Message: "Prefer `command -v` over `type`. `type` output is not stable/standard for checking command existence.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "command -v usage",
			input:    `command -v ls`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1064")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1065(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "normal command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "properly spaced brackets",
			input:    `[ -f /tmp/foo ]`,
			expected: []katas.Violation{},
		},
		{
			name:     "properly spaced double brackets",
			input:    `[[ -f /tmp/foo ]]`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1065")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1066(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid loop over command",
			input:    `for i in $(ls); do echo $i; done`,
			expected: []katas.Violation{}, // ZC1066 targets cat only, though ZC1050 might catch ls
		},
		{
			name:     "valid loop over glob",
			input:    `for i in *; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid loop over cat",
			input: `for i in $(cat file); do echo $i; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1066",
					Message: "Avoid iterating over `cat` output. Use `while read` loop or `($(<file))` for line-based iteration.",
					Line:    1,
					Column:  10,
				},
			},
		},
		{
			name:  "invalid loop over backtick cat",
			input: "for i in `cat file`; do echo $i; done",
			expected: []katas.Violation{
				{
					KataID:  "ZC1066",
					Message: "Avoid iterating over `cat` output. Use `while read` loop or `($(<file))` for line-based iteration.",
					Line:    1,
					Column:  10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1066")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1067(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid export separation",
			input:    `var=$(cmd); export var`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid export literal",
			input:    `export var="value"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid export command substitution",
			input: `export var=$(cmd)`,
			expected: []katas.Violation{
				{
					KataID: "ZC1067",
					Message: "Exporting and assigning a command substitution in one step masks the return value. " +
						"Use `var=$(cmd); export var`.",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name:  "invalid export backticks",
			input: "export var=`cmd`",
			expected: []katas.Violation{
				{
					KataID: "ZC1067",
					Message: "Exporting and assigning a command substitution in one step masks the return value. " +
						"Use `var=$(cmd); export var`.",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name:     "valid export with no assignment",
			input:    `export var`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1067")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1068(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid hook registration",
			input:    `autoload -Uz add-zsh-hook; add-zsh-hook precmd my_precmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid normal function",
			input:    `my_func() { echo hello; }`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid precmd definition",
			input: `precmd() { echo "prompt"; }`,
			expected: []katas.Violation{
				{
					KataID: "ZC1068",
					Message: "Defining `precmd` directly overwrites existing hooks. " +
						"Use `autoload -Uz add-zsh-hook; add-zsh-hook precmd my_func` instead.",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name:  "invalid chpwd definition",
			input: `function chpwd() { ls; }`,
			expected: []katas.Violation{
				{
					KataID: "ZC1068",
					Message: "Defining `chpwd` directly overwrites existing hooks. " +
						"Use `autoload -Uz add-zsh-hook; add-zsh-hook chpwd my_func` instead.",
					Line: 1,
					// Start of "function" keyword usually, or name depending on parser.
					Column: 1,
				},
			},
		},
		{
			// Regression for #1225 — anonymous function must not panic.
			name:     "anonymous function is safe",
			input:    `() { echo hi } "$@"`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1068")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1069(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		// ZC1069 is inert. In Zsh `local` equals `typeset` and is valid at
		// any scope, so none of these inputs may report a violation.
		{
			name:     "local in function",
			input:    `my_func() { local x=1; }`,
			expected: []katas.Violation{},
		},
		{
			name:     "typeset global",
			input:    `typeset x=1`,
			expected: []katas.Violation{},
		},
		{
			name:     "local at global scope (valid in zsh)",
			input:    `local x=1`,
			expected: []katas.Violation{},
		},
		{
			name:     "local in if block (valid in zsh)",
			input:    `if true; then local x=1; fi`,
			expected: []katas.Violation{},
		},
		{
			name:     "local in nested function",
			input:    `outer() { inner() { local x=1; }; }`,
			expected: []katas.Violation{},
		},
		{
			name:     "local in subshell (valid in zsh)",
			input:    `( local x=1 )`,
			expected: []katas.Violation{},
		},
		{
			name:     "local in function keyword",
			input:    "function myfunc { local x=1; }",
			expected: []katas.Violation{},
		},
		{
			name:     "local in while loop (valid in zsh)",
			input:    `while true; do local x=1; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "local in for loop (valid in zsh)",
			input:    `for i in a b c; do local x=1; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "local in case (valid in zsh)",
			input:    "case $x in\na) local y=1;;\nesac",
			expected: []katas.Violation{},
		},
		{
			name:     "regular echo command",
			input:    `echo hello world`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1069")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1070(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid builtin wrapper",
			input:    `cd() { builtin cd "$@"; }`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid command wrapper",
			input:    `ls() { command ls --color "$@"; }`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid recursive wrapper",
			input: `cd() { cd "$@"; }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1070",
					Message: "Recursive call to `cd` inside `cd`. Use `builtin cd` or `command cd` to invoke the underlying command.",
					Line:    1,
					Column:  8, // Position of inner `cd`
				},
			},
		},
		{
			name:  "invalid recursive ls wrapper",
			input: `ls() { ls -G; }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1070",
					Message: "Recursive call to `ls` inside `ls`. Use `builtin ls` or `command ls` to invoke the underlying command.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:     "valid recursive custom function (ignored)",
			input:    `myfunc() { echo hi; myfunc; }`,
			expected: []katas.Violation{},
		},
		{
			name: "valid recursive with condition (false positive risk)",
			// Static analysis warns anyway because it sees direct recursion.
			// ZC1070 intends to warn about WRAPPERS where you usually mean builtin.
			// For algorithms, recursion is valid.
			// Maybe limit ZC1070 to common builtins?
			// Or just warn "Recursive call ... ensure this is intended or use builtin".
			// The message says "Use builtin ...".
			// If it's an algorithm, `builtin myfunc` is invalid (unless myfunc is a builtin?).
			// `builtin` only works for builtins. `command` works for external.
			// If `myfunc` is a function, `command myfunc` ignores function? Yes.
			// So if I want to call the *function* recursively, I use `myfunc`.
			// So ZC1070 flagged valid recursion as error if I imply "infinite recursion".

			// Let's stick to checking builtins to be safe?
			// Or accept that "recursive function" warning is useful but wording should change.
			// "Recursive call detected. If wrapping a builtin/command, use `builtin` or `command`."
			// But for standard recursion `fib(n-1)`, this warning is annoying.

			// Decision: Limit to known builtins + `ls`, `grep` etc?
			// Or just "standard builtins".
			// Let's update logic to only flag if name is in a "common wrapper targets" list?
			// Or "common builtins".

			input:    `fib() { fib $(($1-1)); }`,
			expected: []katas.Violation{}, // Should NOT warn for generic recursion?
		},
		{
			// Regression for #1226 — anonymous function (nil Name) must
			// not panic even when the check walks its body.
			name:     "anonymous function is safe",
			input:    `() { cd /tmp } "$@"`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1070")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1071(t *testing.T) {
	// t.Skip("Skipping ZC1071 tests due to parser limitation with array literals. See issue #41.")

	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid append self reference single",
			input: `arr=($arr)`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1071",
					Message: "Appending to an array using `arr=($arr ...)` is verbose and slower. Use `arr+=(...)` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "valid append whole array first then new elements",
			input: `arr=($arr a b)`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1071",
					Message: "Appending to an array using `arr=($arr ...)` is verbose and slower. Use `arr+=(...)` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "valid append whole array via bracket-at form",
			input: `arr=(${arr[@]} x)`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1071",
					Message: "Appending to an array using `arr=($arr ...)` is verbose and slower. Use `arr+=(...)` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			// Prepend: the new element comes first, so `+=` would append
			// to the wrong end. `path=($new $path)` is not a true append.
			name:     "prepend is not an append",
			input:    `path=($new $path)`,
			expected: []katas.Violation{},
		},
		{
			// Transform: `${files##*/}` rewrites every element; `+=` is
			// meaningless here.
			name:     "transform is not an append",
			input:    `files=(${files##*/})`,
			expected: []katas.Violation{},
		},
		{
			// Truncate: `${tokens[0,-2]}` drops the last element; `+=`
			// cannot express that.
			name:     "truncate is not an append",
			input:    `tokens=(${tokens[0,-2]})`,
			expected: []katas.Violation{},
		},
		{
			name:     "no self reference in array",
			input:    `arr=(a b c)`,
			expected: []katas.Violation{},
		},
		{
			name:     "simple assignment no array",
			input:    `x=5`,
			expected: []katas.Violation{},
		},
		{
			name:     "non-assignment operator",
			input:    `x + 5;`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1071")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1072(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid awk",
			input:    `awk '/pattern/ {print}' file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid grep recursive",
			input:    `grep -r pattern . | awk '{print $1}'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep awk",
			input: `grep pattern file | awk '{print $1}'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1072",
					Message: "Use `awk '/pattern/ {...}'` instead of `grep pattern | awk '{...}'` to avoid a pipeline.",
					Line:    1,
					Column:  19, // Position of pipe
				},
			},
		},
		{
			name:  "invalid grep awk with flags",
			input: `grep -i pattern file | awk '{print}'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1072",
					Message: "Use `awk '/pattern/ {...}'` instead of `grep pattern | awk '{...}'` to avoid a pipeline.",
					Line:    1,
					Column:  22,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1072")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1073(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid arithmetic",
			input:    `(( i = i + 1 ))`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid special vars",
			input:    `(( $# > 0 ))`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid positional",
			input:    `(( $1 > 5 ))`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid expansion",
			input:    `(( ${#arr} > 0 ))`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid simple variable",
			input: `(( $i > 5 ))`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1073",
					Message: "Unnecessary use of `$` in arithmetic expressions. Use `(( var ))` instead of `(( $var ))`.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:  "invalid multiple",
			input: `(( $x + $y ))`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1073",
					Message: "Unnecessary use of `$` in arithmetic expressions. Use `(( var ))` instead of `(( $var ))`.",
					Line:    1,
					Column:  4,
				},
				{
					KataID:  "ZC1073",
					Message: "Unnecessary use of `$` in arithmetic expressions. Use `(( var ))` instead of `(( $var ))`.",
					Line:    1,
					Column:  9,
				},
			},
		},
		{
			name:     "valid command subst",
			input:    `(( $(date +%s) > 0 ))`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1073")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1074(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid modifier usage",
			input:    `echo ${path:h} ${file:t}`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid dirname usage",
			input: `dir=$(dirname $path)`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1074",
					Message: "Use '${var:h}' instead of '$(dirname $var)'. Modifiers are faster and built-in.",
					Line:    1,
					Column:  5,
				},
			},
		},
		{
			name:  "invalid basename usage",
			input: `base=$(basename $path)`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1074",
					Message: "Use '${var:t}' instead of '$(basename $var)'. Modifiers are faster and built-in.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			name:  "invalid backtick dirname",
			input: "dir=`dirname $path`",
			expected: []katas.Violation{
				{
					KataID:  "ZC1074",
					Message: "Use '${var:h}' instead of '$(dirname $var)'. Modifiers are faster and built-in.",
					Line:    1,
					Column:  5,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1074")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1075(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "quoted variable",
			input:    `rm "$var"`,
			expected: []katas.Violation{},
		},
		{
			name:  "unquoted variable",
			input: `rm $var`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1075",
					Message: "Quote `$var`. An unquoted empty or unset value is elided entirely, dropping the word.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:  "unquoted array access",
			input: `ls ${files[1]}`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1075",
					Message: "Quote this array element. An unquoted empty value is elided, dropping the word.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:     "quoted array access",
			input:    `ls "${files[1]}"`,
			expected: []katas.Violation{},
		},
		{
			// A suffixed expansion keeps a literal tail (`/file`), so it
			// cannot elide and must not be flagged.
			name:     "suffixed expansion does not elide",
			input:    `cp $src/file dest`,
			expected: []katas.Violation{},
		},
		{
			// The RHS of an assignment word given to an assignment
			// builtin is not glob/split-expanded in Zsh.
			name:     "export assignment word is not glob-subject",
			input:    `export FOO=$BAR`,
			expected: []katas.Violation{},
		},
		{
			name:     "readonly assignment word is not glob-subject",
			input:    `readonly R=$S`,
			expected: []katas.Violation{},
		},
		{
			name:     "integer assignment word is not glob-subject",
			input:    `integer n=$m`,
			expected: []katas.Violation{},
		},
		{
			// A bare name argument (no `=`) to export is a normal
			// argument and remains flagged.
			name:  "export of a dynamic name is still flagged",
			input: `export $dynamic`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1075",
					Message: "Quote `$dynamic`. An unquoted empty or unset value is elided entirely, dropping the word.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			// `:=` default provides a value, so the expansion never
			// elides — the canonical `: ${VAR:=default}` idiom.
			name:     "assign-default modifier does not elide",
			input:    `: ${VAR:=default}`,
			expected: []katas.Violation{},
		},
		{
			name:     "minus-default modifier does not elide",
			input:    `echo ${V:-fallback}`,
			expected: []katas.Violation{},
		},
		{
			name:     "plus-alternate modifier is not flagged",
			input:    `echo ${V:+set}`,
			expected: []katas.Violation{},
		},
		{
			// A path modifier (`:h`) can still produce an empty word, so
			// it stays flagged — and a scalar is not mislabelled as an
			// array element.
			name:  "path modifier scalar is flagged as expansion",
			input: `echo ${V:h}`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1075",
					Message: "Quote this expansion. An unquoted empty or unset value is elided, dropping the word.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1075")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1076(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid autoload",
			input:    `autoload -Uz my_func`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid autoload split flags",
			input:    `autoload -U -z my_func`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid autoload with other flags",
			input:    `autoload -UzX my_func`,
			expected: []katas.Violation{},
		},
		{
			name:  "missing U",
			input: `autoload -z my_func`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1076",
					Message: "Use `autoload -Uz` to ensure consistent and safe function loading.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "missing z",
			input: `autoload -U my_func`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1076",
					Message: "Use `autoload -Uz` to ensure consistent and safe function loading.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "missing flags",
			input: `autoload my_func`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1076",
					Message: "Use `autoload -Uz` to ensure consistent and safe function loading.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1076")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1077(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid expansion",
			input:    `upper=${var:u}`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid tr upper",
			input: `upper=$(echo $var | tr 'a-z' 'A-Z')`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1077",
					Message: "Use `${var:u}` instead of `tr` for uppercase conversion. It is faster and built-in.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:  "invalid tr lower",
			input: `lower=$(echo $var | tr 'A-Z' 'a-z')`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1077",
					Message: "Use `${var:l}` instead of `tr` for lowercase conversion. It is faster and built-in.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:  "invalid tr upper posix",
			input: `upper=$(echo $var | tr '[:lower:]' '[:upper:]')`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1077",
					Message: "Use `${var:u}` instead of `tr` for uppercase conversion. It is faster and built-in.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:     "other tr usage",
			input:    `clean=$(echo $var | tr -d '\n')`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1077")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1078(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "quoted arguments",
			input:    `cmd "$@"`,
			expected: []katas.Violation{},
		},
		{
			name:     "quoted star",
			input:    `cmd "$*"`,
			expected: []katas.Violation{},
		},
		{
			name:  "unquoted arguments",
			input: `cmd $@`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1078",
					Message: "Unquoted $@ drops empty elements. Use \"$@\" to preserve every positional parameter.",
					Line:    1,
					Column:  5,
				},
			},
		},
		{
			name:  "unquoted star",
			input: `cmd $*`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1078",
					Message: "Unquoted $* drops empty elements. Use \"$*\" to preserve every positional parameter.",
					Line:    1,
					Column:  5,
				},
			},
		},
		{
			name:  "mixed",
			input: `cmd arg1 $@ arg2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1078",
					Message: "Unquoted $@ drops empty elements. Use \"$@\" to preserve every positional parameter.",
					Line:    1,
					Column:  10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1078")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1079(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid quoted comparison",
			input:    `[[ $var == "$other" ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid literal comparison",
			input:    `[[ $var == "foo" ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid pattern comparison (literal)",
			input:    `[[ $var == foo* ]]`,
			expected: []katas.Violation{},
		},
		{
			// In default Zsh the RHS variable is matched literally (no
			// GLOB_SUBST), so quoting it is a no-op and nothing fires.
			name:     "unquoted variable RHS is literal in zsh",
			input:    `[[ $var == $other ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "unquoted variable RHS with != is literal",
			input:    `[[ $var != $other ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "unquoted array access RHS is literal",
			input:    `[[ $var = ${arr[1]} ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "non-equality operator",
			input:    `[[ $a -lt $b ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "quoted string RHS",
			input:    `[[ $a == "hello" ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "integer comparison",
			input:    `[[ $a -eq 5 ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "no infix elements",
			input:    `[[ -f /tmp/foo ]]`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1079")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1080(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid nullglob",
			input:    `for f in *.txt(N); do echo $f; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid no glob",
			input:    `for f in a b c; do echo $f; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid variable",
			input:    `for f in $files; do echo $f; done`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid glob star",
			input: `for f in *.txt; do echo $f; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1080",
					Message: "Glob '*.txt' will error if no matches found. Append `(N)` to make it nullglob.",
					Line:    1,
					Column:  10,
				},
			},
		},
		{
			name:  "invalid glob question",
			input: `for f in file?; do echo $f; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1080",
					Message: "Glob 'file?' will error if no matches found. Append `(N)` to make it nullglob.",
					Line:    1,
					Column:  10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1080")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1081(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid length check",
			input:    `len=${#var}`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid echo pipe wc -c",
			input: `len=$(echo $var | wc -c)`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1081",
					Message: "Use `${#var}` to get string length. Pipeline to `wc` is inefficient.",
					Line:    1,
					Column:  17,
				},
			},
		},
		{
			name:  "invalid print pipe wc -m",
			input: `len=$(print -r $var | wc -m)`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1081",
					Message: "Use `${#var}` to get string length. Pipeline to `wc` is inefficient.",
					Line:    1,
					Column:  21,
				},
			},
		},
		{
			name:     "wc on file (valid)",
			input:    `wc -c file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "cat pipe wc (valid-ish)",
			input:    `cat file | wc -c`,
			expected: []katas.Violation{}, // ZC1038 might flag cat usage, but ZC1081 shouldn't
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1081")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1082(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid expansion",
			input:    `new=${var//foo/bar}`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sed s///",
			input: `new=$(echo $var | sed 's/foo/bar/')`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1082",
					Message: "Use `${var//old/new}` for string replacement. Pipeline to `sed` is inefficient.",
					Line:    1,
					Column:  17,
				},
			},
		},
		{
			name:  "invalid sed s///g",
			input: `new=$(echo $var | sed "s/foo/bar/g")`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1082",
					Message: "Use `${var//old/new}` for string replacement. Pipeline to `sed` is inefficient.",
					Line:    1,
					Column:  17,
				},
			},
		},
		{
			name:  "invalid sed different separator",
			input: `new=$(print $var | sed 's|foo|bar|')`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1082",
					Message: "Use `${var//old/new}` for string replacement. Pipeline to `sed` is inefficient.",
					Line:    1,
					Column:  18,
				},
			},
		},
		{
			name:     "valid sed other usage",
			input:    `echo $var | sed -n '/p/p'`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1082")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1083(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid brace expansion",
			input:    `echo {1..10}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid brace expansion with list",
			input:    `echo {a,b,c}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid variable after braces",
			input:    `echo {1..10}$var`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid variable before braces",
			input:    `echo $var{1..10}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid variable inside list expansion",
			input:    `echo {a,b,$var}`,
			expected: []katas.Violation{},
		},
		{
			// Zsh resolves `$n` to an integer before brace expansion, so
			// `{1..$n}` expands correctly. This is the Bash-only failure
			// the kata used to flag; under Zsh it is valid.
			name:     "valid variable as range end in Zsh",
			input:    `echo {1..$n}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid variable as range start in Zsh",
			input:    `echo {$n..10}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid variable as range start and end in Zsh",
			input:    `echo {$min..$max}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid command substitution in range in Zsh",
			input:    `echo {1..$(echo 10)}`,
			expected: []katas.Violation{},
		},
		{
			// A param substitution that merely contains `{`, `..`, and `$`
			// is not a brace range — the `{` is preceded by `$`.
			name:     "valid param substitution that looks like a range",
			input:    `echo "${branch/foo../bar}"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid param substitution replace empty",
			input:    `echo "${x/../}"`,
			expected: []katas.Violation{},
		},
		{
			// A quoted brace range with a parameter bound: the quotes
			// suppress the expansion that the unquoted form gets in Zsh,
			// so the value stays the literal string. A `{` preceded by `$`
			// (a `${…}` expansion) is excluded by zc1083HasBareBraceOpen.
			name:  "invalid quoted brace range stays literal",
			input: `echo "{1..$n}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1083",
					Message: "A quoted brace range stays literal. `\"{1..$n}\"` does not expand; in Zsh the unquoted `{1..$n}` does. Drop the quotes to expand the range.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1083")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1084(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "quoted glob",
			input:    `find . -name "*.txt"`,
			expected: []katas.Violation{},
		},
		{
			name:     "single quoted glob",
			input:    `find . -name '*.txt'`,
			expected: []katas.Violation{},
		},
		{
			name:  "unquoted star glob",
			input: `find . -name *.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1084",
					Message: "Quote globs in `find` commands. `*.txt` is subject to shell expansion.",
					Line:    1,
					Column:  14,
				},
			},
		},
		{
			name:  "unquoted question glob",
			input: `find . -name file?.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1084",
					Message: "Quote globs in `find` commands. `file?.txt` is subject to shell expansion.",
					Line:    1,
					Column:  14,
				},
			},
		},
		{
			name:  "unquoted bracket glob (merged)",
			input: `find . -name[a-z]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1084",
					Message: "Quote globs in `find` commands. `-(name[a-z])` contains unquoted brackets.",
					Line:    1,
					Column:  13, // Points to [
				},
			},
		},
		{
			name:  "unquoted bracket glob (space)",
			input: `find . -name [a-z]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1084",
					Message: "Quote globs in `find` commands. `[a-z]` is subject to shell expansion.",
					Line:    1,
					Column:  14, // Points to [
				},
			},
		},
		{
			name:  "unquoted bracket glob (partial)",
			input: `find . -name file[a-z]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1084",
					Message: "Quote globs in `find` commands. `file[a-z]` is subject to shell expansion.",
					Line:    1,
					Column:  18, // Points to [
				},
			},
		},
		{
			name:     "escaped glob",
			input:    `find . -name \*.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "escaped question",
			input:    `find . -name file\?.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "double backslash (escaped backslash + glob)",
			input: `find . -name \\*.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1084",
					Message: "Quote globs in `find` commands. `\\\\*.txt` is subject to shell expansion.",
					Line:    1,
					Column:  14,
				},
			},
		},

		{
			name:     "other flag (ignore)",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:     "find with exec (ignore args)",
			input:    `find . -exec echo {} \;`,
			expected: []katas.Violation{},
		},
		{
			name:     "quoted bracket glob",
			input:    `find . -name '[a-z]'`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vioations := testutil.Check(tt.input, "ZC1084")
			testutil.AssertViolations(t, tt.input, vioations, tt.expected)
		})
	}
}

func TestZC1085(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		// Unquoted array iteration is the correct Zsh idiom — quoting
		// would collapse the array into one word — so the rule no
		// longer warns on any of these forms.
		{
			name:     "unquoted array iteration is correct and not flagged",
			input:    `for i in $items; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "unquoted array access not flagged",
			input:    `for i in ${items[@]}; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "unquoted command substitution not flagged",
			input:    `for i in $(ls); do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "mixed literal and unquoted not flagged",
			input:    `for i in start $items end; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "quoted forms remain clean",
			input:    `for i in "${items[@]}"; do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "glob expansion not flagged",
			input:    `for i in *.txt; do echo $i; done`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1085")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1086(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid function definition",
			input:    `my_func() { echo hello; }`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid function definition with keyword",
			input: `function my_func { echo hello; }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1086",
					Message: "Prefer `func() { ... }` over `function func { ... }` for portability and consistency.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid function definition with keyword and parens",
			input: `function my_func() { echo hello; }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1086",
					Message: "Prefer `func() { ... }` over `function func { ... }` for portability and consistency.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1086")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1087(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid redirection to new file",
			input:    `cat input.txt > output.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid append redirection",
			input:    `cat file.txt >> file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid clobbering simple command",
			input: `sort file.txt > file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1087",
					Message: "Output redirection overwrites input file `file.txt`. The file is truncated before reading.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid clobbering pipeline",
			input: `cat file.txt | grep foo > file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1087",
					Message: "Output redirection overwrites input file `file.txt`. The file is truncated before reading.",
					Line:    1,
					Column:  14, // Points to |
				},
			},
		},
		{
			name:  "invalid clobbering with input redirection",
			input: `grep foo < file.txt > file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1087",
					Message: "Output redirection overwrites input file `file.txt`. The file is truncated before reading.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid different files",
			input:    `sed 's/a/b/' input > output`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1087")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1088(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid subshell",
			input:    `( ls )`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid subshell with return",
			input:    `( return 1 )`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid subshell with exit",
			input:    `( exit 1 )`,
			expected: []katas.Violation{},
		},
		/*
			{
				name:     "valid subshell checked exit status",
				input:    `( cd /tmp ) || exit`,
				expected: []katas.Violation{},
			},
		*/
		{
			name:     "valid subshell used in condition",
			input:    `if ( cd /tmp ); then :; fi`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid subshell side effect",
			input: `( cd /tmp )`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1088",
					Message: "Subshell `( ... )` isolates state changes. The changes (e.g. `cd`, variable assignment) will be lost. Use `{ ... }` to preserve them, or add commands that use the changed state.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid subshell variable assignment",
			input: `( var=1 )`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1088",
					Message: "Subshell `( ... )` isolates state changes. The changes (e.g. `cd`, variable assignment) will be lost. Use `{ ... }` to preserve them, or add commands that use the changed state.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid subshell output capture",
			input:    `out=$( ( cd /tmp; pwd ) )`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid subshell export",
			input: `( export VAR=1 )`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1088",
					Message: "Subshell `( ... )` isolates state changes. The changes (e.g. `cd`, variable assignment) will be lost. Use `{ ... }` to preserve them, or add commands that use the changed state.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "brace group is fine",
			input:    `{ cd /tmp; }`,
			expected: []katas.Violation{},
		},
		{
			name:     "no subshell just command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "subshell with non-state-changing commands",
			input:    `( echo hello; ls )`,
			expected: []katas.Violation{},
		},
		{
			name:     "while loop no violation",
			input:    `while true; do echo x; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "for loop no violation",
			input:    `for i in a b c; do echo $i; done`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1088")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1089(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid redirection order",
			input:    `cmd > file 2>&1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid ampersand redirection",
			input:    `cmd &> file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid redirection order",
			input: `cmd 2>&1 > file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1089",
					Message: "Redirection order matters. `2>&1 > file` does not redirect stderr to file. Use `> file 2>&1` instead.",
					Line:    1,
					Column:  10, // Points to > (outer)
				},
			},
		},
		{
			name:  "invalid redirection order append",
			input: `cmd 2>&1 >> file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1089",
					Message: "Redirection order matters. `2>&1 > file` does not redirect stderr to file. Use `> file 2>&1` instead.",
					Line:    1,
					Column:  11,
				},
			},
		},
		{
			name:     "unrelated redirection",
			input:    `cmd 2>&3 > file`,
			expected: []katas.Violation{},
		},
		{
			name:  "redirection to file named 1",
			input: `cmd >& 1 > file`, // >& 1 means to file named 1 if 1 is IDENT? But lexer makes 1 INT.
			// In zsh `cmd >& 1` is same as `cmd &> 1`?
			// `2>&1` is explicit.
			// If I write `cmd >&1`, it is `cmd` `>&` `1`.
			// If checkZC1089 is strict about `2` arg...
			// My check doesn't check left side of inner redirection.
			// It assumes `... >& 1`.
			// If `cmd >& 1 > file`.
			// Inner `cmd >& 1`.
			// Is `cmd` stderr/stdout?
			// `>&` redirects stdout AND stderr.
			// So `(cmd >& 1) > file`.
			// stdout+stderr -> 1.
			// Then stdout (of result?) -> file.
			// Since stdout was redirected to 1, result stdout is empty?
			// So `> file` is empty?
			// This is also weird but not the specific `2>&1` mistake.
			// I should strictly check for `2` on the left of inner `>&`?
			// But my parser puts `2` as ARGUMENT to command.
			// `SimpleCommand` args: `cmd`, `2`.
			// `Redirection` Left is `SimpleCommand`.
			// How to check if `2` is the last argument?
			// `redir.Left` -> SimpleCommand. Check last arg is "2".
			// But wait, if `cmd arg 2>&1`.
			// `SimpleCommand` args: `cmd`, `arg`, `2`.
			// So verify last arg is "2".
			expected: []katas.Violation{}, // Should ideally update test to be correct or check logic
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1089")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1090(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid unquoted regex",
			input:    `[[ $v =~ ^foo ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid mixed regex",
			input:    `[[ $v =~ "user_"* ]]`,
			expected: []katas.Violation{},
		},
		// In Zsh the =~ RHS is a regex regardless of quoting (unlike
		// Bash), and quoting protects metacharacters from globbing and
		// word splitting — so none of these quoted forms are flagged.
		{
			name:     "quoted start anchor is a valid Zsh regex",
			input:    `[[ $v =~ "^foo" ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "quoted wildcard is a valid Zsh regex",
			input:    `[[ $v =~ "foo.*" ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "quoted alternation is a valid Zsh regex",
			input:    `[[ $v =~ "a|b" ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "quoted regex with a space needs its quotes",
			input:    `[[ $v =~ "^a b$" ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "quoted literal",
			input:    `[[ $v =~ "foo" ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "quoted variable with meta",
			input:    `[[ $v =~ "^$pat" ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "non-regex operator ignored",
			input:    `[[ $a == "^foo" ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "no double bracket",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1090")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1091(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid string comparison",
			input:    `[[ $a == $b ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid file test",
			input:    `[[ -f $file ]]`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid arithmetic -eq",
			input: `[[ $a -eq 1 ]]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1091",
					Message: "Use `(( ... ))` for arithmetic comparisons. e.g. `(( a < b ))` instead of `[[ a -lt b ]]`.",
					Line:    1,
					Column:  7, // -eq token column
				},
			},
		},
		{
			name:  "invalid arithmetic -lt",
			input: `[[ $a -lt 5 ]]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1091",
					Message: "Use `(( ... ))` for arithmetic comparisons. e.g. `(( a < b ))` instead of `[[ a -lt b ]]`.",
					Line:    1,
					Column:  7,
				},
			},
		},
		{
			name:  "invalid nested arithmetic",
			input: `[[ $a -gt 0 && $b -lt 10 ]]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1091",
					Message: "Use `(( ... ))` for arithmetic comparisons. e.g. `(( a < b ))` instead of `[[ a -lt b ]]`.",
					Line:    1,
					Column:  7, // -gt
				},
				{
					KataID:  "ZC1091",
					Message: "Use `(( ... ))` for arithmetic comparisons. e.g. `(( a < b ))` instead of `[[ a -lt b ]]`.",
					Line:    1,
					Column:  19, // -lt
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1091")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1092(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			// ZC1092's generic "prefer print" advice is consolidated
			// under ZC1037 (`print -r --`), so it no longer fires.
			name:     "echo defers to ZC1037",
			input:    `echo "hello world"`,
			expected: []katas.Violation{},
		},
		{
			name:     "echo with flags defers to ZC1037",
			input:    `echo -n "hello"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid print",
			input:    `print "hello world"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid printf",
			input:    `printf "%s\n" "hello"`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1092")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

// ZC1093 was retired as a duplicate of ZC1038. Kept as a no-op stub so
// legacy `disabled_katas` lists keep parsing; the detection runs under
// ZC1038 now.

func TestZC1093Stub(t *testing.T) {
	cases := []string{
		"echo hi",
		"ls",
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			v := testutil.Check(in, "ZC1093")
			testutil.AssertViolations(t, in, v, []katas.Violation{})
		})
	}
}

func TestZC1094(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sed with -i flag",
			input:    `sed -i 's/foo/bar/' file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sed with file argument",
			input:    `sed 's/foo/bar/' file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sed with -e flag",
			input:    `sed -e 's/foo/bar/' -e 's/baz/qux/'`,
			expected: []katas.Violation{},
		},
		{
			name:     "bare sed reads a stream not a variable",
			input:    `sed 's/foo/bar/'`,
			expected: []katas.Violation{},
		},
		{
			name:     "sed below a non-echo pipe has no variable to expand",
			input:    `git branch | sed 's/foo/bar/'`,
			expected: []katas.Violation{},
		},
		{
			name:  "echo of a variable into sed",
			input: `echo $v | sed 's/foo/bar/'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1094",
					Message: "Use `${var//pattern/replacement}` instead of piping through `sed` for simple substitutions. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  11,
				},
			},
		},
		{
			name:  "here-string of a variable into sed",
			input: `sed 's/foo/bar/' <<< $v`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1094",
					Message: "Use `${var//pattern/replacement}` instead of piping through `sed` for simple substitutions. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1094")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1095(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid seq with range",
			input:    `seq 1 10`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid seq with step",
			input:    `seq 1 2 10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid seq with single number",
			input: `seq 5`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1095",
					Message: "Use `repeat N do ... done` or `for i in {1..N}` instead of `seq N`. Zsh has built-in constructs for repetition that avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid seq with large number",
			input: `seq 100`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1095",
					Message: "Use `repeat N do ... done` or `for i in {1..N}` instead of `seq N`. Zsh has built-in constructs for repetition that avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid non-numeric argument",
			input:    `seq abc`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1095")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1096(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "bc usage",
			input: `echo "1.5 + 2.5" | bc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1096",
					Message: "Zsh supports floating point arithmetic natively. You often don't need `bc`.",
					Line:    1,
					Column:  20,
				},
			},
		},
		{
			name:     "valid arithmetic",
			input:    `(( 1.5 + 2.5 ))`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1096")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1097(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid local loop variable",
			input:    `my_func() { local i; for i in 1 2; do echo $i; done; }`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid global loop variable",
			input: `my_func() { for i in 1 2; do echo $i; done; }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1097",
					Message: "Loop variable 'i' is used without 'local'. It will be global. Use `local i` before the loop.",
					Line:    1,
					Column:  17,
				},
			},
		},
		{
			name:  "invalid global loop variable (implicit in)",
			input: `my_func() { for i; do echo $i; done; }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1097",
					Message: "Loop variable 'i' is used without 'local'. It will be global. Use `local i` before the loop.",
					Line:    1,
					Column:  17,
				},
			},
		},
		{
			name:     "valid typeset loop variable",
			input:    `my_func() { typeset i; for i in 1 2; do echo $i; done; }`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid integer loop variable",
			input:    `my_func() { integer i; for i in 1 2; do echo $i; done; }`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid local loop variable in nested block",
			input:    `my_func() { if true; then local i; for i in 1 2; do echo $i; done; fi; }`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid local inside loop (too late)",
			input: `my_func() { for i in 1 2; do local i; echo $i; done; }`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1097",
					Message: "Loop variable 'i' is used without 'local'. It will be global. Use `local i` before the loop.",
					Line:    1,
					Column:  17,
				},
			},
		},
		{
			name:     "valid arithmetic for loop (C-style)",
			input:    `my_func() { for ((i=0; i<10; i++)); do echo $i; done; }`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1097")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1098(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "no eval",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "eval without variables",
			input:    `eval "echo hello"`,
			expected: []katas.Violation{},
		},
		{
			name:     "eval of a bare variable command string",
			input:    `eval $cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "eval of a quoted standalone variable",
			input:    `eval "$REPLY"`,
			expected: []katas.Violation{},
		},
		{
			name:     "eval of a single flagged expansion",
			input:    `eval ${(e)x}`,
			expected: []katas.Violation{},
		},
		{
			name:  "eval with unquoted variable",
			input: `eval "ls $dir"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1098",
					Message: "Use the `(q)` flag (or `(qq)`, `(q-)`) when using variables in `eval` to prevent injection.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "eval with quoted variable",
			input:    `eval "ls ${(q)dir}"`,
			expected: []katas.Violation{},
		},
		{
			// `(q)` cannot quote a command substitution, so the
			// shell-init idiom must not flag.
			name:     "eval with command substitution only",
			input:    `eval "$(pyenv init - zsh)"`,
			expected: []katas.Violation{},
		},
		{
			name:  "eval mixing command substitution and parameter",
			input: `eval "$(cmd) ${userdata}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1098",
					Message: "Use the `(q)` flag (or `(qq)`, `(q-)`) when using variables in `eval` to prevent injection.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1098")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1099(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "while read loop with pipe",
			input: `cat file | while read line; do echo $line; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1099",
					Message: "Consider using `for line in ${(f)variable}` instead of `... | while read line`. It's faster and cleaner in Zsh.",
					Line:    1,
					Column:  10,
				},
			},
		},
		{
			name:     "while loop without pipe",
			input:    `while read line; do echo $line; done < file`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1099")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
