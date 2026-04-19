package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1906(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt POSIX_CD` (explicit default)",
			input:    `unsetopt POSIX_CD`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt AUTO_PUSHD` (unrelated)",
			input:    `setopt AUTO_PUSHD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt POSIX_CD`",
			input: `setopt POSIX_CD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1906",
					Message: "`setopt POSIX_CD` changes when `cd`/`pushd` read `CDPATH` — scripts that relied on Zsh's default silently enter different directories. Keep it off; wrap POSIX-specific code with `emulate -LR sh`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_POSIX_CD`",
			input: `unsetopt NO_POSIX_CD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1906",
					Message: "`unsetopt NO_POSIX_CD` changes when `cd`/`pushd` read `CDPATH` — scripts that relied on Zsh's default silently enter different directories. Keep it off; wrap POSIX-specific code with `emulate -LR sh`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1906")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
