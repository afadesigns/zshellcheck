package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1993(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt KSH_TYPESET` (default)",
			input:    `unsetopt KSH_TYPESET`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_KSH_TYPESET`",
			input:    `setopt NO_KSH_TYPESET`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt KSH_TYPESET`",
			input: `setopt KSH_TYPESET`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1993",
					Message: "`setopt KSH_TYPESET` re-splits the RHS of every later `typeset`/`local` — `typeset path=$HOME/My Files` now treats `Files` as a second name. Scope with `emulate -LR ksh` inside the one helper that needs it.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_KSH_TYPESET`",
			input: `unsetopt NO_KSH_TYPESET`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1993",
					Message: "`unsetopt NO_KSH_TYPESET` re-splits the RHS of every later `typeset`/`local` — `typeset path=$HOME/My Files` now treats `Files` as a second name. Scope with `emulate -LR ksh` inside the one helper that needs it.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1993")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
