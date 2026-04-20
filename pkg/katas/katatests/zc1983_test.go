package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1983(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt CSH_JUNKIE_QUOTES`",
			input:    `unsetopt CSH_JUNKIE_QUOTES`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_CSH_JUNKIE_QUOTES`",
			input:    `setopt NO_CSH_JUNKIE_QUOTES`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt CSH_JUNKIE_QUOTES`",
			input: `setopt CSH_JUNKIE_QUOTES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1983",
					Message: "`setopt CSH_JUNKIE_QUOTES` makes every later multi-line `\"…\"`/`'…'` an error — inlined SQL/JSON payloads and autoloaded helpers stop parsing. Scope csh-style strictness with `emulate -LR csh` in the one helper that needs it.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_CSH_JUNKIE_QUOTES`",
			input: `unsetopt NO_CSH_JUNKIE_QUOTES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1983",
					Message: "`unsetopt NO_CSH_JUNKIE_QUOTES` makes every later multi-line `\"…\"`/`'…'` an error — inlined SQL/JSON payloads and autoloaded helpers stop parsing. Scope csh-style strictness with `emulate -LR csh` in the one helper that needs it.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1983")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
