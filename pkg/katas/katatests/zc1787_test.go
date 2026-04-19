package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1787(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt EXTENDED_GLOB` (unrelated)",
			input:    `setopt EXTENDED_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `set -e` (unrelated)",
			input:    `set -e`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt AUTO_CD`",
			input: `setopt AUTO_CD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1787",
					Message: "`setopt AUTO_CD` turns any bare directory name into a silent `cd`. A typo or a user-controlled value reshapes `$PWD`; keep this in `~/.zshrc`, not in scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt autocd` (lowercase / no underscore)",
			input: `setopt autocd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1787",
					Message: "`setopt autocd` turns any bare directory name into a silent `cd`. A typo or a user-controlled value reshapes `$PWD`; keep this in `~/.zshrc`, not in scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `set -o autocd`",
			input: `set -o autocd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1787",
					Message: "`set -o autocd` turns any bare directory name into a silent `cd`. A typo or a user-controlled value reshapes `$PWD`; keep this in `~/.zshrc`, not in scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1787")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
