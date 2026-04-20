package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1977(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt CHASE_DOTS`",
			input:    `unsetopt CHASE_DOTS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_CHASE_DOTS`",
			input:    `setopt NO_CHASE_DOTS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt CHASE_DOTS`",
			input: `setopt CHASE_DOTS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1977",
					Message: "`setopt CHASE_DOTS` makes `cd ..` physically resolve before walking up — blue/green `current` symlinks stop working for `../foo` lookups. Keep off; use `cd -P` one-shot when physical resolution is needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_CHASE_DOTS`",
			input: `unsetopt NO_CHASE_DOTS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1977",
					Message: "`unsetopt NO_CHASE_DOTS` makes `cd ..` physically resolve before walking up — blue/green `current` symlinks stop working for `../foo` lookups. Keep off; use `cd -P` one-shot when physical resolution is needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1977")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
