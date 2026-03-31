package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1244(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid mv -n",
			input:    `mv -n src dst`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid mv -f explicit",
			input:    `mv -f old new`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid bare mv",
			input: `mv file.txt backup.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1244",
					Message: "Consider `mv -n` to prevent overwriting existing files. Without `-n`, `mv` silently overwrites the target.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1244")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
