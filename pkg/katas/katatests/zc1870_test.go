package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1870(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt GLOB_ASSIGN` (explicit default)",
			input:    `unsetopt GLOB_ASSIGN`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt GLOB_ASSIGN`",
			input: `setopt GLOB_ASSIGN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1870",
					Message: "`setopt GLOB_ASSIGN` expands glob patterns on the RHS of `var=` — `logs=*.log` silently captures the first match, `cert=~/secrets/*` picks up attacker drops. Keep it off; use explicit `arr=( *.log )`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_GLOB_ASSIGN`",
			input: `unsetopt NO_GLOB_ASSIGN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1870",
					Message: "`unsetopt NO_GLOB_ASSIGN` expands glob patterns on the RHS of `var=` — `logs=*.log` silently captures the first match, `cert=~/secrets/*` picks up attacker drops. Keep it off; use explicit `arr=( *.log )`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1870")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
