package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1150(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid cat with flags",
			input:    `cat -n file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid cat multiple files",
			input:    `cat file1 file2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat single file",
			input: `cat config.yaml`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1150",
					Message: "Use `$(<file)` instead of `$(cat file)` to read file contents. Zsh's `$(<file)` is a built-in that avoids spawning cat.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1150")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
