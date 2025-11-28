package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1080(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid nullglob",
			input:    `for f in *.txt(N); do echo $f; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid no glob",
			input:    `for f in a b c; do echo $f; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid variable",
			input:    `for f in $files; do echo $f; done`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid glob star",
			input: `for f in *.txt; do echo $f; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1080",
					Message: "Glob '*.txt' will error if no matches found. Append `(N)` to make it nullglob.",
					Line:    1,
					Column:  10,
				},
			},
		},
		{
			name:  "invalid glob question",
			input: `for f in file?; do echo $f; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1080",
					Message: "Glob 'file?' will error if no matches found. Append `(N)` to make it nullglob.",
					Line:    1,
					Column:  10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1080")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
