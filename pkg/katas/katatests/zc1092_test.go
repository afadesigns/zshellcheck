package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1092(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid echo",
			input: `echo "hello world"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1092",
					Message: "Prefer `print` over `echo`. `echo` behavior varies. `print` is the Zsh builtin. Especially with flags, `print -n` or `print -r` is more reliable.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid echo with flags",
			input: `echo -n "hello"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1092",
					Message: "Prefer `print` over `echo`. `echo` behavior varies. `print` is the Zsh builtin. Especially with flags, `print -n` or `print -r` is more reliable.",
					Line:    1,
					Column:  1,
				},
			},
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
