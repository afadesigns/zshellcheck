package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1936(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt POSIX_ALIASES` (explicit default)",
			input:    `unsetopt POSIX_ALIASES`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EXTENDED_GLOB` (unrelated)",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt POSIX_ALIASES`",
			input: `setopt POSIX_ALIASES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1936",
					Message: "`setopt POSIX_ALIASES` narrows alias expansion to plain identifiers — aliases on `if`/`for`/`function` silently stop firing and any library that hooked them breaks. Scope with `emulate -LR sh` instead of flipping globally.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_POSIX_ALIASES`",
			input: `unsetopt NO_POSIX_ALIASES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1936",
					Message: "`unsetopt NO_POSIX_ALIASES` narrows alias expansion to plain identifiers — aliases on `if`/`for`/`function` silently stop firing and any library that hooked them breaks. Scope with `emulate -LR sh` instead of flipping globally.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1936")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
