package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1859(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt MULTIOS` (explicit default)",
			input:    `setopt MULTIOS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt MULTIOS`",
			input: `unsetopt MULTIOS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1859",
					Message: "`unsetopt MULTIOS` reverts to POSIX single-output redirection — `cmd >a >b` silently drops `a`, log collectors stop receiving new lines. Keep the option on; scope inside a `LOCAL_OPTIONS` function if one line really needs POSIX.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_MULTIOS`",
			input: `setopt NO_MULTIOS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1859",
					Message: "`setopt NO_MULTIOS` reverts to POSIX single-output redirection — `cmd >a >b` silently drops `a`, log collectors stop receiving new lines. Keep the option on; scope inside a `LOCAL_OPTIONS` function if one line really needs POSIX.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1859")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
