package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1110(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid head with file",
			input:    `head -1 file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid head -5",
			input:    `head -5`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid tail with file",
			input:    `tail -1 file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid head -1 in pipeline",
			input: `head -1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1110",
					Message: "Use `${lines[1]}` instead of `head -1`. Zsh array subscripts avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid tail -1 in pipeline",
			input: `tail -1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1110",
					Message: "Use `${lines[-1]}` instead of `tail -1`. Zsh array subscripts avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid head -n 1",
			input: `head -n 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1110",
					Message: "Use `${lines[1]}` instead of `head -1`. Zsh array subscripts avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1110")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
