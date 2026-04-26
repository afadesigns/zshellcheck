// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1300(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ZSH_VERSION usage",
			input:    `echo $ZSH_VERSION`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_VERSINFO usage",
			input: `echo $BASH_VERSINFO`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1300",
					Message: "Avoid Bash version variables in Zsh — use `$ZSH_VERSION` instead. Bash version variables are undefined in Zsh.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			name:  "invalid BASH_VERSION usage",
			input: `echo $BASH_VERSION`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1300",
					Message: "Avoid Bash version variables in Zsh — use `$ZSH_VERSION` instead. Bash version variables are undefined in Zsh.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1300")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1301(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid pipestatus usage",
			input:    `echo $pipestatus`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid PIPESTATUS usage",
			input: `echo $PIPESTATUS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1301",
					Message: "Avoid `$PIPESTATUS` in Zsh — use `$pipestatus` (lowercase) instead. The uppercase form is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1301")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1302(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid man usage",
			input:    `man zshbuiltins`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid help usage",
			input: `help cd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1302",
					Message: "Avoid `help` in Zsh — it is a Bash builtin. Use `run-help` or `man zshbuiltins` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1302")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1303(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid zmodload usage",
			input:    `zmodload zsh/stat`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid enable without -f",
			input:    `enable mybuiltin`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid enable -f usage",
			input: `enable -f /path/to/builtin mybuiltin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1303",
					Message: "Avoid `enable -f` in Zsh — use `zmodload` to load modules. `enable -f` is Bash-specific.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1303")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1304(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ZSH_SUBSHELL usage",
			input:    `echo $ZSH_SUBSHELL`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_SUBSHELL usage",
			input: `echo $BASH_SUBSHELL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1304",
					Message: "Avoid `$BASH_SUBSHELL` in Zsh — use `$ZSH_SUBSHELL` instead. `BASH_SUBSHELL` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1304")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1305(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid words usage",
			input:    `echo $words`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid COMP_WORDS usage",
			input: `echo $COMP_WORDS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1305",
					Message: "Avoid `$COMP_WORDS` in Zsh — use `$words` array instead. `COMP_WORDS` is Bash completion-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1305")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1306(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid CURRENT usage",
			input:    `echo $CURRENT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid COMP_CWORD usage",
			input: `echo $COMP_CWORD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1306",
					Message: "Avoid `$COMP_CWORD` in Zsh — use `$CURRENT` instead. `COMP_CWORD` is Bash completion-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1306")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1307(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid dirstack usage",
			input:    `echo $dirstack`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid DIRSTACK usage",
			input: `echo $DIRSTACK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1307",
					Message: "Avoid `$DIRSTACK` in Zsh — use `$dirstack` (lowercase) instead. The uppercase form is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1307")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1308(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid BUFFER usage",
			input:    `echo $BUFFER`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid COMP_LINE usage",
			input: `echo $COMP_LINE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1308",
					Message: "Avoid `$COMP_LINE` in Zsh — use `$BUFFER` instead. `COMP_LINE` is Bash completion-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1308")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1309(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid non-Bash variable",
			input:    `echo $MY_COMMAND`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_COMMAND usage",
			input: `echo $BASH_COMMAND`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1309",
					Message: "Avoid `$BASH_COMMAND` in Zsh — it is undefined. Use `$ZSH_DEBUG_CMD` in debug traps if needed.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1309")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1310(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid variable",
			input:    `echo $MY_STRING`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_EXECUTION_STRING usage",
			input: `echo $BASH_EXECUTION_STRING`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1310",
					Message: "Avoid `$BASH_EXECUTION_STRING` in Zsh — it is undefined. Access command arguments directly instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1310")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1311(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid compdef usage",
			input:    `compdef _git git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid complete usage",
			input: `complete -F _git git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1311",
					Message: "Avoid `complete` in Zsh — it is a Bash builtin. Use `compdef` for Zsh completion registration.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1311")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1312(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid compadd usage",
			input:    `compadd foo bar baz`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid compgen usage",
			input: `compgen -W "foo bar" -`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1312",
					Message: "Avoid `compgen` in Zsh — it is a Bash builtin. Use `compadd` or Zsh completion functions instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1312")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1313(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid aliases usage",
			input:    `echo $aliases`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_ALIASES usage",
			input: `echo $BASH_ALIASES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1313",
					Message: "Avoid `$BASH_ALIASES` in Zsh — use the `aliases` associative array instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1313")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1314(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid non-Bash variable",
			input:    `echo $MY_PATH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_LOADABLES_PATH",
			input: `echo $BASH_LOADABLES_PATH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1314",
					Message: "Avoid `$BASH_LOADABLES_PATH` in Zsh — it is undefined. Use `zmodload` with full module names.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1314")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1315(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid emulate usage",
			input:    `emulate -L sh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_COMPAT usage",
			input: `echo $BASH_COMPAT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1315",
					Message: "Avoid `$BASH_COMPAT` in Zsh — use `emulate` for shell compatibility mode instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1315")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1316(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid funcfiletrace usage",
			input:    `echo $funcfiletrace`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid caller usage",
			input: `caller 0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1316",
					Message: "Avoid `caller` in Zsh — it is a Bash builtin. Use `$funcfiletrace` and `$funcstack` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1316")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1317(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ZDOTDIR usage",
			input:    `echo $ZDOTDIR`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_ENV usage",
			input: `echo $BASH_ENV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1317",
					Message: "Avoid `$BASH_ENV` in Zsh — use `$ZDOTDIR` for Zsh startup file locations instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1317")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1318(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid commands usage",
			input:    `echo $commands`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_CMDS usage",
			input: `echo $BASH_CMDS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1318",
					Message: "Avoid `$BASH_CMDS` in Zsh — use the `$commands` hash for command path lookups instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1318")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1319(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid argument count",
			input:    `echo $MYVAR`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_ARGC usage",
			input: `echo $BASH_ARGC`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1319",
					Message: "Avoid `$BASH_ARGC` in Zsh — use `$#` for argument count. `BASH_ARGC` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1319")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1320(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid argv usage",
			input:    `echo $argv`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_ARGV usage",
			input: `echo $BASH_ARGV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1320",
					Message: "Avoid `$BASH_ARGV` in Zsh — use `$argv` or `$@` for positional parameters. `BASH_ARGV` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1320")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1321(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid variable",
			input:    `echo $MY_FD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_XTRACEFD usage",
			input: `echo $BASH_XTRACEFD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1321",
					Message: "Avoid `$BASH_XTRACEFD` in Zsh — it is undefined. Redirect stderr directly for xtrace output.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1321")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1322(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid non-COPROC variable",
			input:    `echo $MY_PROC`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid COPROC usage",
			input: `echo $COPROC`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1322",
					Message: "Avoid `$COPROC` in Zsh — Zsh coprocesses use `read -p`/`print -p` for I/O. `COPROC` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1322")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1323(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid kill usage",
			input:    `kill -STOP 1234`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid suspend usage",
			input: `suspend -f`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1323",
					Message: "Avoid `suspend` in Zsh — it is a Bash builtin. Use `kill -STOP $$` if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1323")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1324(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid precmd usage",
			input:    `echo $precmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid PROMPT_COMMAND usage",
			input: `echo $PROMPT_COMMAND`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1324",
					Message: "Avoid `$PROMPT_COMMAND` in Zsh — use the `precmd` hook function instead. `PROMPT_COMMAND` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1324")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1325(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid PS1 usage",
			input:    `echo $PS1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid PS0 usage",
			input: `echo $PS0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1325",
					Message: "Avoid `$PS0` in Zsh — use the `preexec` hook function instead. `PS0` is Bash 4.4+ specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1325")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1326(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid HISTFILE usage",
			input:    `echo $HISTFILE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid HISTTIMEFORMAT usage",
			input: `echo $HISTTIMEFORMAT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1326",
					Message: "Avoid `$HISTTIMEFORMAT` in Zsh — use `setopt EXTENDED_HISTORY` and `fc -li` instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1326")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1327(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid fc usage",
			input:    `fc -l`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid history without flags",
			input:    `history 10`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `history -c` is owned by ZC1487",
			input:    `history -c`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `history -w` (Bash-only write)",
			input: `history -w`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1327",
					Message: "Avoid `history -w` in Zsh — Bash history flags differ. Use `fc` commands for Zsh history management.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `history -a`",
			input: `history -a`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1327",
					Message: "Avoid `history -a` in Zsh — Bash history flags differ. Use `fc` commands for Zsh history management.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1327")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1328(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid HISTSIZE usage",
			input:    `echo $HISTSIZE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid HISTCONTROL usage",
			input: `echo $HISTCONTROL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1328",
					Message: "Avoid `$HISTCONTROL` in Zsh — use `setopt HIST_IGNORE_DUPS` and related options instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1328")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1329(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid HISTSIZE usage",
			input:    `echo $HISTSIZE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid HISTIGNORE usage",
			input: `echo $HISTIGNORE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1329",
					Message: "Avoid `$HISTIGNORE` in Zsh — use `zshaddhistory` hook for history filtering instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1329")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1330(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ZDOTDIR usage",
			input:    `echo $ZDOTDIR`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid INPUTRC usage",
			input: `echo $INPUTRC`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1330",
					Message: "Avoid `$INPUTRC` in Zsh — Zsh uses `bindkey` and ZLE, not readline. `INPUTRC` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1330")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1331(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid match usage",
			input:    `echo $match`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_REMATCH usage",
			input: `echo $BASH_REMATCH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1331",
					Message: "Avoid `$BASH_REMATCH` in Zsh — use `$match` array and `$MATCH` for regex captures instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1331")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1332(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid variable",
			input:    `echo $MYGLOB`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid GLOBIGNORE usage",
			input: `echo $GLOBIGNORE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1332",
					Message: "Avoid `$GLOBIGNORE` in Zsh — use `setopt EXTENDED_GLOB` with `~` operator for glob exclusion.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1332")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1333(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid TIMEFMT usage",
			input:    `echo $TIMEFMT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid TIMEFORMAT usage",
			input: `echo $TIMEFORMAT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1333",
					Message: "Avoid `$TIMEFORMAT` in Zsh — use `$TIMEFMT` instead. Format specifiers differ between Bash and Zsh.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1333")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1334(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid whence -p usage",
			input:    `whence -p git`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid type without -p",
			input:    `type git`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid type -p usage",
			input: `type -p git`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1334",
					Message: "Avoid `type -p` in Zsh — use `whence -p` to get the command path instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1334")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1335(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid other command",
			input:    `cat file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "tac usage",
			input: `tac file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1335",
					Message: "Consider Zsh `${(Oa)array}` for reversing array data instead of piping to `tac`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1335")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1336(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid typeset -x usage",
			input:    `typeset -x PATH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid printenv usage",
			input: `printenv HOME`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1336",
					Message: "Avoid `printenv` in Zsh — use `typeset -x` or `export` to list environment variables.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1336")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1337(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid print usage",
			input:    `print -l hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "fold usage",
			input: `fold -w 80 file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1337",
					Message: "Consider Zsh `$COLUMNS` and `print` for text wrapping instead of `fold`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1337")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1338(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid seq without -s",
			input:    `seq 10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid seq -s usage",
			input: `seq -s , 10`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1338",
					Message: "Avoid `seq -s` in Zsh — use `${(j:sep:)array}` with brace expansion for joined sequences.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1338")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1339(t *testing.T) {
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
			name:     "valid wc -c",
			input:    `wc -c`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid wc -l in pipeline",
			input: `wc -l`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1339",
					Message: "Use Zsh `${#${(f)var}}` for line counting instead of piping through `wc -l`. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1339")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1340(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — non-shuf command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — shuf -n 1",
			input: `shuf -n 1 file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1340",
					Message: "Avoid `shuf` for random selection — use Zsh `${array[RANDOM%$#array+1]}` with `$RANDOM` for in-shell randomness without spawning an external.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1340")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1341(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find with other predicate",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -executable",
			input: `find . -executable`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1341",
					Message: "Use Zsh `*(.x)` glob qualifier instead of `find -executable`. The `.` restricts to regular files and `x` to the executable bit.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -executable with -type f",
			input: `find . -type f -executable`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1341",
					Message: "Use Zsh `*(.x)` glob qualifier instead of `find -executable`. The `.` restricts to regular files and `x` to the executable bit.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1341")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1342(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -empty",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -empty",
			input: `find . -empty`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1342",
					Message: "Use Zsh `*(L0)` glob qualifier instead of `find -empty`. Add `.` for regular files only: `*(.L0)`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -type f -empty",
			input: `find . -type f -empty -delete`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1342",
					Message: "Use Zsh `*(L0)` glob qualifier instead of `find -empty`. Add `.` for regular files only: `*(.L0)`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1342")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1343(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without age predicate",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -mtime +7",
			input: `find . -mtime +7`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1343",
					Message: "Use Zsh glob qualifiers (`*(m±N)`, `*(M±N)`, `*(a±N)`, `*(c±N)`) instead of `find -mtime`/`-mmin`/`-atime`/`-amin`/`-ctime`/`-cmin` for age predicates.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -mmin -60",
			input: `find . -mmin -60`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1343",
					Message: "Use Zsh glob qualifiers (`*(m±N)`, `*(M±N)`, `*(a±N)`, `*(c±N)`) instead of `find -mtime`/`-mmin`/`-atime`/`-amin`/`-ctime`/`-cmin` for age predicates.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -ctime 0",
			input: `find . -ctime 0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1343",
					Message: "Use Zsh glob qualifiers (`*(m±N)`, `*(M±N)`, `*(a±N)`, `*(c±N)`) instead of `find -mtime`/`-mmin`/`-atime`/`-amin`/`-ctime`/`-cmin` for age predicates.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1343")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1344(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -size",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -size +10M",
			input: `find . -size +10M`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1344",
					Message: "Use Zsh `*(L±N[kmp])` glob qualifier instead of `find -size`. Unit suffixes: `k` kilobytes, `m` megabytes, `p` 512-byte blocks.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -size -1k",
			input: `find . -size -1k`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1344",
					Message: "Use Zsh `*(L±N[kmp])` glob qualifier instead of `find -size`. Unit suffixes: `k` kilobytes, `m` megabytes, `p` 512-byte blocks.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1344")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1345(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -perm",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -perm 755",
			input: `find . -perm 755`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1345",
					Message: "Use Zsh `*(f:mode:)` glob qualifier instead of `find -perm`. Octal (`*(f:0755:)`) or symbolic (`*(f:u+x:)`) expressions are both supported.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -perm -u+x",
			input: `find . -perm -u+x`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1345",
					Message: "Use Zsh `*(f:mode:)` glob qualifier instead of `find -perm`. Octal (`*(f:0755:)`) or symbolic (`*(f:u+x:)`) expressions are both supported.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1345")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1346(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -user",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -user",
			input: `find . -user alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1346",
					Message: "Use Zsh `*(u:name:)` / `*(u+uid)` / `*(U)` glob qualifiers instead of `find -user`/`-uid`/`-nouser`. Ownership predicates live entirely in the shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -uid",
			input: `find / -uid 1000`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1346",
					Message: "Use Zsh `*(u:name:)` / `*(u+uid)` / `*(U)` glob qualifiers instead of `find -user`/`-uid`/`-nouser`. Ownership predicates live entirely in the shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -nouser",
			input: `find / -nouser`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1346",
					Message: "Use Zsh `*(u:name:)` / `*(u+uid)` / `*(U)` glob qualifiers instead of `find -user`/`-uid`/`-nouser`. Ownership predicates live entirely in the shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1346")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1347(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -group",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -group",
			input: `find . -group wheel`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1347",
					Message: "Use Zsh `*(g:name:)` / `*(g+gid)` / `*(G)` glob qualifiers instead of `find -group`/`-gid`/`-nogroup`. Group predicates live entirely in the shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -gid 10",
			input: `find . -gid 10`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1347",
					Message: "Use Zsh `*(g:name:)` / `*(g+gid)` / `*(G)` glob qualifiers instead of `find -group`/`-gid`/`-nogroup`. Group predicates live entirely in the shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1347")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1348(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -type",
			input:    `find . -name "*.txt"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -type f",
			input: `find . -type f`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1348",
					Message: "Use Zsh glob type qualifiers (`*(/)`, `*(.)`, `*(@)`, `*(=)`, `*(p)`, `*(*)`, `*(%)`) instead of `find -type`. No external process required.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -type d",
			input: `find / -type d`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1348",
					Message: "Use Zsh glob type qualifiers (`*(/)`, `*(.)`, `*(@)`, `*(=)`, `*(p)`, `*(*)`, `*(%)`) instead of `find -type`. No external process required.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1348")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1349(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — expr arithmetic (not length)",
			input:    `expr 1 + 2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — expr length with var",
			input: `expr length "$s"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1349",
					Message: "Use `${#var}` instead of `expr length \"$var\"` for string length. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — expr length with literal",
			input: `expr length hello`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1349",
					Message: "Use `${#var}` instead of `expr length \"$var\"` for string length. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1349")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1350(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — expr arithmetic (not substr)",
			input:    `expr 2 + 3`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — expr substr",
			input: `expr substr "$s" 1 5`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1350",
					Message: "Use `${str:pos:len}` instead of `expr substr` for substring extraction. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1350")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1351(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — expr arithmetic",
			input:    `expr 1 + 2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — expr match",
			input: `expr match "$s" '^foo'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1351",
					Message: "Use `[[ $str =~ pattern ]]` with `$match` / `$MATCH` arrays instead of `expr match`/`expr index`. Regex evaluation stays in the shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — expr index",
			input: `expr index "$s" aeiou`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1351",
					Message: "Use `[[ $str =~ pattern ]]` with `$match` / `$MATCH` arrays instead of `expr match`/`expr index`. Regex evaluation stays in the shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1351")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1352(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — xargs without -I",
			input:    `xargs -n 1 echo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — xargs -I{}",
			input: `xargs -I{} echo hi`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1352",
					Message: "Avoid `xargs -I{}` — iterate with `for x in ${(f)\"$(cmd)\"}; do ...; done` in Zsh. A for loop avoids one-subprocess-per-item and keeps variables in scope.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — xargs -Ixx custom replace-string",
			input: `xargs -Ifile cp file /tmp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1352",
					Message: "Avoid `xargs -I{}` — iterate with `for x in ${(f)\"$(cmd)\"}; do ...; done` in Zsh. A for loop avoids one-subprocess-per-item and keeps variables in scope.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1352")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1353(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — printf without -v",
			input:    `printf 'hello %s\n' world`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — printf -v var",
			input: `printf -v line '%d' 42`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1353",
					Message: "Avoid `printf -v` in Zsh — use `print -v var -rf fmt ...` or `var=$(printf fmt ...)`. `-v` is Bash-specific and ignored elsewhere.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1353")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1354(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — type without flags",
			input:    `type grep`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — type -t",
			input: `type -t grep`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1354",
					Message: "Use Zsh `whence -w` (category), `whence -a` (all), or `whence -p` (path) instead of Bash-specific `type -t`/`-a`/`-P`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — type -P",
			input: `type -P grep`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1354",
					Message: "Use Zsh `whence -w` (category), `whence -a` (all), or `whence -p` (path) instead of Bash-specific `type -t`/`-a`/`-P`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1354")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1355(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo without -E",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo -E raw",
			input: `echo -E "$line"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1355",
					Message: "Use `print -r` instead of `echo -E` for raw output. `-E` is a Bash-ism and ignored by POSIX echo.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1355")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1356(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — read -A",
			input:    `read -A arr`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — read -r line",
			input:    `read -r line`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — read -a (Bash syntax)",
			input: `read -a arr`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1356",
					Message: "Use `read -A` (uppercase) in Zsh to read into an array. `read -a` has different semantics in Zsh than in Bash.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1356")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1357(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — printf %s",
			input:    `printf '%s\n' "$line"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — printf %q",
			input: `printf '%q' "$v"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1357",
					Message: "Use Zsh `${(q)var}` for shell-quoting instead of `printf '%q'`. Variants: `${(qq)}`, `${(qqq)}`, `${(qqqq)}` for different quote styles.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1357")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1358(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pwd without -P",
			input:    `pwd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pwd -P",
			input: `pwd -P`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1358",
					Message: "Use `${PWD:P}` instead of `pwd -P` — the `P` modifier resolves symlinks and returns the canonical path without spawning an external.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1358")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1359(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — id without group flag",
			input:    `id -u`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — id -Gn",
			input: `id -Gn`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1359",
					Message: "Avoid `id -Gn`/`-G`/`-gn`/`-g` — use Zsh `$groups` (names→gids assoc array) or `$GID` for the primary group after `zmodload zsh/parameter`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — id -g",
			input: `id -g`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1359",
					Message: "Avoid `id -Gn`/`-G`/`-gn`/`-g` — use Zsh `$groups` (names→gids assoc array) or `$GID` for the primary group after `zmodload zsh/parameter`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1359")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1360(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ls -l (not sort-by-size)",
			input:    `ls -l`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ls -S",
			input: `ls -S`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1360",
					Message: "Use Zsh `*(OL)` (largest-first) or `*(oL)` (smallest-first) glob qualifier instead of `ls -S`. No external process needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ls -lS",
			input: `ls -lS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1360",
					Message: "Use Zsh `*(OL)` (largest-first) or `*(oL)` (smallest-first) glob qualifier instead of `ls -S`. No external process needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1360")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1361(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — awk with generic program",
			input:    `awk '{print $1}' file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — awk NR==5",
			input: `awk 'NR==5' file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1361",
					Message: "Avoid `awk 'NR==N'` — split with `${(f)\"$(<file)\"}` in Zsh and index: `lines=(${(f)\"$(<file)\"}); print $lines[N]`. No awk process needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1361")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1362(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — test without -o",
			input:    `test -f file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — test -o noglob",
			input: `test -o noglob`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1362",
					Message: "Use `[[ -o option ]]` for option checks in Zsh — `test -o` means logical OR, not option-query, producing wrong results.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1362")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1363(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -newer",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -newer",
			input: `find . -newer ref.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1363",
					Message: "Use Zsh `*(e:'[[ $REPLY -nt REF ]]':)` eval glob qualifier instead of `find -newer`/`-anewer`/`-cnewer`/`-newerXY`. `$REPLY` holds the current match.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1363")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1364(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — cut -f (field, different kata)",
			input:    `cut -f 2 file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — cut -c",
			input: `cut -c 1-5 file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1364",
					Message: "Use Zsh `${var:pos:len}` for character ranges instead of `cut -c`. Parameter expansion is in-shell and zero-indexed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — cut -c attached",
			input: `cut -c1-5 file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1364",
					Message: "Use Zsh `${var:pos:len}` for character ranges instead of `cut -c`. Parameter expansion is in-shell and zero-indexed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1364")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1365(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — stat without format flag",
			input:    `stat file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — stat -c %s",
			input: `stat -c %s file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1365",
					Message: "Use Zsh `zmodload zsh/stat; zstat -H meta file` for file metadata instead of `stat -c '%...'`. The associative array `meta` exposes every stat field.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — stat -c %Y (mtime)",
			input: `stat -c %Y file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1365",
					Message: "Use Zsh `zmodload zsh/stat; zstat -H meta file` for file metadata instead of `stat -c '%...'`. The associative array `meta` exposes every stat field.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1365")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1366(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — limit builtin",
			input:    `limit cputime 10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ulimit",
			input: `ulimit -t 10`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1366",
					Message: "Use Zsh `limit` (human-readable) or `limit -s` (stdout-only) instead of POSIX `ulimit` for Zsh-native resource queries.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1366")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1367(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — printf normal format",
			input:    `printf '%s\n' hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — printf %(fmt)T",
			input: `printf '%(%Y-%m-%d)T\n' 1700000000`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1367",
					Message: "Use Zsh `strftime fmt seconds` (from `zsh/datetime`) instead of Bash `printf '%(fmt)T' seconds`. Same formatting, more readable, no Bash-version gating.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1367")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1368(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sh without -c",
			input:    `sh script.sh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sh -c",
			input: `sh -c 'echo hi'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1368",
					Message: "Avoid `sh -c`/`bash -c` inside a Zsh script — inline the code as a function to keep access to arrays, associative arrays, and Zsh features. Use `zsh -c` only when a fresh shell is truly required.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — bash -c",
			input: `bash -c 'echo hi'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1368",
					Message: "Avoid `sh -c`/`bash -c` inside a Zsh script — inline the code as a function to keep access to arrays, associative arrays, and Zsh features. Use `zsh -c` only when a fresh shell is truly required.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1368")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1369(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — od -x (hex, different use)",
			input:    `od -x file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — od -c",
			input: `od -c file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1369",
					Message: "Use Zsh `${(V)var}` to see non-printable characters in a variable — renders control chars as `\\n`, `\\t`, etc., without spawning `od`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1369")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1370(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — non-yes command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — yes str",
			input: `yes banana`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1370",
					Message: "Prefer Zsh `repeat N { print str }` over `yes str | head -n N` for producing N copies of a line. No external `yes` process, no pipe.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1370")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1371(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — basename single path",
			input:    `basename /usr/bin/zsh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — basename -a",
			input: `basename -a /a/b /c/d`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1371",
					Message: "Use Zsh `${paths:t}` on an array for bulk basename extraction instead of `basename -a`. The `:t` modifier applies to every array element.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1371")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1372(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — zmv usage",
			input:    `zmv '(*).txt' '$1.md'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — rename perl-style",
			input: `rename 's/\.txt$/.md/' *.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1372",
					Message: "Use Zsh `zmv` (autoload -Uz zmv) instead of `rename`/`rename.ul`/`prename`. Glob-pattern renaming is handled in-shell with capture groups.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — rename.ul util-linux",
			input: `rename.ul .txt .md *.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1372",
					Message: "Use Zsh `zmv` (autoload -Uz zmv) instead of `rename`/`rename.ul`/`prename`. Glob-pattern renaming is handled in-shell with capture groups.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1372")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1373(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — env without -0",
			input:    `env VAR=val cmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — env -0",
			input: `env -0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1373",
					Message: "Use Zsh `${(0)\"$(<file)\"}` to split NUL-terminated content in-shell. `env -0` is usually followed by `xargs -0` or a read loop — both avoided.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1373")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1374(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo other var",
			input:    `echo $var`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $FUNCNEST expecting depth",
			input: `echo $FUNCNEST`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1374",
					Message: "In Zsh, `$FUNCNEST` is the configured limit, not the current depth. Use `${#funcstack}` for current function nesting depth.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1374")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1375(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — tty (print tty name)",
			input:    `tty`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tty -s",
			input: `tty -s`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1375",
					Message: "Use `[[ -t 0 ]]` (stdin), `[[ -t 1 ]]` (stdout), or `[[ -t 2 ]]` (stderr) instead of `tty -s`. In-shell file-descriptor test, no external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1375")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1376(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated echo",
			input:    `echo $VAR`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $BASH_XTRACEFD",
			input: `echo $BASH_XTRACEFD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1376",
					Message: "`BASH_XTRACEFD` is Bash-only. Zsh ignores it. Redirect trace output with `exec {fd}>file; exec 2>&$fd; setopt XTRACE` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1376")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1377(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $aliases",
			input:    `echo $aliases`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $BASH_ALIASES",
			input: `echo $BASH_ALIASES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1377",
					Message: "`$BASH_ALIASES` is Bash-only. In Zsh use `$aliases` (assoc array) — same structure, e.g. `print -l ${(kv)aliases}`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1377")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1378(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $dirstack (lowercase)",
			input:    `echo $dirstack`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $DIRSTACK",
			input: `echo $DIRSTACK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1378",
					Message: "Use lowercase `$dirstack` in Zsh — uppercase `$DIRSTACK` is Bash-only.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1378")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1379(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated echo",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $PROMPT_COMMAND",
			input: `echo $PROMPT_COMMAND`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1379",
					Message: "`PROMPT_COMMAND` is Bash-only. In Zsh define a `precmd` function or use `autoload -Uz add-zsh-hook; add-zsh-hook precmd my_hook`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1379")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1380(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $HISTORY_IGNORE (Zsh)",
			input:    `echo $HISTORY_IGNORE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $HISTIGNORE",
			input: `echo $HISTIGNORE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1380",
					Message: "`$HISTIGNORE` is Bash-only. In Zsh use `$HISTORY_IGNORE` (underscored) for the same history-pattern filter.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1380")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1381(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $words (Zsh compsys)",
			input:    `echo $words`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $COMP_WORDS",
			input: `echo $COMP_WORDS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1381",
					Message: "Bash `$COMP_*` completion variables do not exist in Zsh. Use `$words` (array of tokens), `$CURRENT` (cursor index), `$BUFFER`, or the `_arguments`/`_values` helpers from `compsys`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — echo $COMP_CWORD",
			input: `echo $COMP_CWORD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1381",
					Message: "Bash `$COMP_*` completion variables do not exist in Zsh. Use `$words` (array of tokens), `$CURRENT` (cursor index), `$BUFFER`, or the `_arguments`/`_values` helpers from `compsys`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1381")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1382(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $BUFFER (Zsh ZLE)",
			input:    `echo $BUFFER`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $READLINE_LINE",
			input: `echo $READLINE_LINE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1382",
					Message: "Bash `$READLINE_*` vars do not exist in Zsh. Inside ZLE widgets use `$BUFFER`, `$CURSOR`, `$MARK`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — echo $READLINE_POINT",
			input: `echo $READLINE_POINT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1382",
					Message: "Bash `$READLINE_*` vars do not exist in Zsh. Inside ZLE widgets use `$BUFFER`, `$CURSOR`, `$MARK`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1382")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1383(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $TIMEFMT (Zsh)",
			input:    `echo $TIMEFMT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $TIMEFORMAT",
			input: `echo $TIMEFORMAT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1383",
					Message: "`$TIMEFORMAT` is Bash-only. Zsh reads `$TIMEFMT` (shorter name) for the `time` builtin's output format.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1383")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1384(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated echo",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $EXECIGNORE",
			input: `echo $EXECIGNORE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1384",
					Message: "`$EXECIGNORE` is Bash-only. For completion filtering in Zsh use `zstyle ':completion:*' ignored-patterns 'pattern'`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1384")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1385(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated echo",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $PS0",
			input: `echo $PS0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1385",
					Message: "`$PS0` is Bash-only. Zsh uses the `preexec` hook function for pre-execution prompts.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1385")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1386(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated echo",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $FIGNORE",
			input: `echo $FIGNORE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1386",
					Message: "`$FIGNORE` is Bash-only. In Zsh use `zstyle ':completion:*' ignored-patterns '*.o *.pyc'` for completion filtering.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1386")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1387(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $options (Zsh)",
			input:    `echo $options`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $SHELLOPTS",
			input: `echo $SHELLOPTS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1387",
					Message: "`$SHELLOPTS` is Bash-only. In Zsh inspect `$options` (assoc array, keys are option names) via `print -l ${(kv)options}`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1387")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1388(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $mailpath (Zsh)",
			input:    `echo $mailpath`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $MAILPATH",
			input: `echo $MAILPATH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1388",
					Message: "Use Zsh lowercase `$mailpath` (array) instead of Bash uppercase `$MAILPATH` (colon-separated string).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1388")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1389(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $hosts (Zsh)",
			input:    `echo $hosts`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $HOSTFILE",
			input: `echo $HOSTFILE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1389",
					Message: "`$HOSTFILE` is Bash-only. Zsh reads hostnames for completion from the `$hosts` array (lowercase).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1389")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1390(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $GROUPS (scalar)",
			input:    `echo $GROUPS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo ${GROUPS[@]}",
			input: `echo ${GROUPS[@]}`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1390",
					Message: "Zsh `$GROUPS` is a scalar (primary GID), not an array. For all group IDs use `${(k)groups}` (after `zmodload zsh/parameter`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1390")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1391(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — test -n VAR",
			input:    `test -n "$VAR"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — test -v VAR",
			input: `test -v VAR`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1391",
					Message: "Use `(( ${+VAR} ))` for Zsh set-check — `-v` is a Bash 4.2+ extension, not reliably portable to Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1391")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1392(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated echo",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $CHILD_MAX",
			input: `echo $CHILD_MAX`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1392",
					Message: "`$CHILD_MAX` is Bash-only. Zsh uses `limit -s maxproc` or `ulimit -u` for process-count limits.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1392")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1393(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $RANDOM",
			input:    `echo $RANDOM`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $SRANDOM",
			input: `echo $SRANDOM`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1393",
					Message: "`$SRANDOM` is Bash 5.1+. In Zsh read `/dev/urandom` directly or use an external (`openssl rand`) for secure random integers.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1393")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1394(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $ZSH_NAME",
			input:    `echo $ZSH_NAME`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $BASH",
			input: `echo $BASH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1394",
					Message: "`$BASH` is Bash-only. Zsh exposes the interpreter name via `$ZSH_NAME` and the executable path indirectly via `$0`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1394")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1395(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — wait $pid",
			input:    `wait $pid`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — wait -n",
			input: `wait -n`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1395",
					Message: "`wait -n` is Bash 4.3+. Zsh's `wait` waits on specific PIDs/jobs or (bare `wait`) all jobs. For any-child semantics, loop over PIDs with individual `wait $pid` calls.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1395")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1396(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unset -v var",
			input:    `unset -v var`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — unset -n ref",
			input: `unset -n ref`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1396",
					Message: "`unset -n` is a Bash nameref operation. Zsh does not honor it; use `unset -v NAME` (variable) or `unset -f NAME` (function) explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1396")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1397(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $compstate (Zsh)",
			input:    `echo $compstate`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $COMP_TYPE",
			input: `echo $COMP_TYPE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1397",
					Message: "Bash `$COMP_TYPE`/`$COMP_KEY`/`$COMP_WORDBREAKS` are not Zsh-native. Use `$compstate` associative array for completion context in Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1397")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1398(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $PS1",
			input:    `echo $PS1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $PROMPT_DIRTRIM",
			input: `echo $PROMPT_DIRTRIM`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1398",
					Message: "`$PROMPT_DIRTRIM` is Bash-only. Use the Zsh prompt escape `%N~` (N = number of path components to keep) for directory truncation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1398")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}

func TestZC1399(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kill signal pid",
			input:    `kill -TERM 1234`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kill -l",
			input: `kill -l`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1399",
					Message: "Use Zsh `print -l $signals` (after `zmodload zsh/parameter`) instead of `kill -l` for listing signal names.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1399")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
