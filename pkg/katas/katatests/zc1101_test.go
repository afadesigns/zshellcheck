package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1101(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid bc with file",
			input:    `bc script.bc`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid bc in pipeline",
			input: `bc -l`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1101",
					Message: "Use `$(( ))` for arithmetic instead of `bc`. Zsh arithmetic expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid other command",
			input:    `calc 1+1`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1101")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
