package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1109(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid cut with file",
			input:    `cut -d: -f1 /etc/passwd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cut with only field",
			input:    `cut -f1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cut in pipeline",
			input: `cut -d: -f1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1109",
					Message: "Use Zsh parameter expansion for field extraction instead of `cut`. `${var%%delim*}` or `${(s.delim.)var}` avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid cut with delimiter and field",
			input: `cut -d',' -f2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1109",
					Message: "Use Zsh parameter expansion for field extraction instead of `cut`. `${var%%delim*}` or `${(s.delim.)var}` avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1109")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
