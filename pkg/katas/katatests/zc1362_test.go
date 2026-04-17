package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1362(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — test without -o",
			input:    `test -f file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — test -o noglob",
			input: `test -o noglob`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1362",
					Message: "Use `[[ -o option ]]` for option checks in Zsh — `test -o` means logical OR, not option-query, producing wrong results.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1362")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
