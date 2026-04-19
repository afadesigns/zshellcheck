package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1916(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt NULL_GLOB` (explicit default)",
			input:    `unsetopt NULL_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EXTENDED_GLOB` (unrelated)",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt NULL_GLOB`",
			input: `setopt NULL_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1916",
					Message: "`setopt NULL_GLOB` makes every later unmatched glob silently empty — `cp *.log /dest` rewrites to `cp /dest`, `rm *.tmp` becomes argv-too-short. Use per-glob `*(N)`, or `setopt LOCAL_OPTIONS NULL_GLOB` in a function.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_NULL_GLOB`",
			input: `unsetopt NO_NULL_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1916",
					Message: "`unsetopt NO_NULL_GLOB` makes every later unmatched glob silently empty — `cp *.log /dest` rewrites to `cp /dest`, `rm *.tmp` becomes argv-too-short. Use per-glob `*(N)`, or `setopt LOCAL_OPTIONS NULL_GLOB` in a function.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1916")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
