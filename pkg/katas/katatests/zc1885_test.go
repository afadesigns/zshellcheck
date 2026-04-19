package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1885(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt CSH_NULL_GLOB` (explicit default)",
			input:    `unsetopt CSH_NULL_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt CSH_NULL_GLOB`",
			input: `setopt CSH_NULL_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1885",
					Message: "`setopt CSH_NULL_GLOB` silently discards unmatched globs in a list when any sibling matches — `rm *.lg *.bak` deletes the `.bak` files and hides the typo. Keep the option off; use `*(N)` per-glob.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_CSH_NULL_GLOB`",
			input: `unsetopt NO_CSH_NULL_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1885",
					Message: "`unsetopt NO_CSH_NULL_GLOB` silently discards unmatched globs in a list when any sibling matches — `rm *.lg *.bak` deletes the `.bak` files and hides the typo. Keep the option off; use `*(N)` per-glob.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1885")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
