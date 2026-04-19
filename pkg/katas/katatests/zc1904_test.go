package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1904(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt KSH_GLOB` (explicit default)",
			input:    `unsetopt KSH_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EXTENDED_GLOB` (unrelated)",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt KSH_GLOB`",
			input: `setopt KSH_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1904",
					Message: "`setopt KSH_GLOB` reinterprets `*(...)` as a ksh-style operator — every Zsh glob qualifier (`*(N)`, `*(D)`, `*(.)`) silently stops working. Prefer `setopt EXTENDED_GLOB`, or scope inside a function with `LOCAL_OPTIONS`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_KSH_GLOB`",
			input: `unsetopt NO_KSH_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1904",
					Message: "`unsetopt NO_KSH_GLOB` reinterprets `*(...)` as a ksh-style operator — every Zsh glob qualifier (`*(N)`, `*(D)`, `*(.)`) silently stops working. Prefer `setopt EXTENDED_GLOB`, or scope inside a function with `LOCAL_OPTIONS`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1904")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
