package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1275(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid Zsh :h modifier",
			input:    `echo ${filepath:h}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid non-dirname command",
			input:    `basename /path/to/file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid dirname usage",
			input: `dirname /path/to/file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1275",
					Message: "Use Zsh parameter expansion `${var:h}` instead of `dirname`. The `:h` modifier extracts the directory without forking a process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1275")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
