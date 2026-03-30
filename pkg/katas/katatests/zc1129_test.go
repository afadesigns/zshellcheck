package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1129(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid wc -l with file",
			input:    `wc -l file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid wc -c without file",
			input:    `wc -c`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid wc -c with file",
			input: `wc -c file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1129",
					Message: "Use `zstat +size file` (via `zmodload zsh/stat`) instead of `wc -c file`. Avoids reading the entire file for a simple size query.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1129")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
