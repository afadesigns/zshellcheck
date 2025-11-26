package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1099(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "while read loop with pipe",
			input: `cat file | while read line; do echo $line; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1099",
					Message: "Consider using `for line in ${(f)variable}` instead of `... | while read line`. It's faster and cleaner in Zsh.",
					Line:    1,
					Column:  10,
				},
			},
		},
		{
			name:  "while loop without pipe",
			input: `while read line; do echo $line; done < file`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1099")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
