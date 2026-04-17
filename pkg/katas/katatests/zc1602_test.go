package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1602(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — setopt NO_NOMATCH",
			input:    `setopt NO_NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — setopt EXTENDED_GLOB",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — setopt KSH_ARRAYS",
			input: `setopt KSH_ARRAYS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1602",
					Message: "`setopt KSH_ARRAYS` flips Zsh core semantics for the whole shell — pre-existing code silently misbehaves. Scope with `emulate -L ksh` / `emulate -L sh` inside the function that needs the mode.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — setopt shwordsplit (lowercase, no underscore)",
			input: `setopt shwordsplit`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1602",
					Message: "`setopt shwordsplit` flips Zsh core semantics for the whole shell — pre-existing code silently misbehaves. Scope with `emulate -L ksh` / `emulate -L sh` inside the function that needs the mode.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1602")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
