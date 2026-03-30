package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1131(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid cat with grep",
			input:    `cat file | grep pattern`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat | read",
			input: `cat file.txt | read line`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1131",
					Message: "Use `while read line; do ...; done < file` instead of `cat file | while read line`. Avoids unnecessary cat and subshell from the pipe.",
					Line:    1,
					Column:  14,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1131")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
