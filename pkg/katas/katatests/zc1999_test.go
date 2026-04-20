package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1999(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt AUTO_NAME_DIRS` (canonical name; handled by ZC1934)",
			input:    `setopt AUTO_NAME_DIRS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt AUTO_NAME_DIRS`",
			input:    `unsetopt AUTO_NAME_DIRS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt AUTO_NAMED_DIRS` (typo)",
			input: `setopt AUTO_NAMED_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1999",
					Message: "`setopt AUTO_NAMED_DIRS` is a typo — the real Zsh option is `AUTO_NAME_DIRS` (no trailing `D`, see ZC1934). Fix the spelling or drop the toggle; `hash -d NAME=PATH` is the explicit alternative.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_AUTO_NAMED_DIRS` (typo)",
			input: `unsetopt NO_AUTO_NAMED_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1999",
					Message: "`unsetopt NO_AUTO_NAMED_DIRS` is a typo — the real Zsh option is `AUTO_NAME_DIRS` (no trailing `D`, see ZC1934). Fix the spelling or drop the toggle; `hash -d NAME=PATH` is the explicit alternative.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1999")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
