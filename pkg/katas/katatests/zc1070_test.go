package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1070")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
