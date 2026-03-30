package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1124(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid cat with file",
			input:    `cat readme.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat /dev/null",
			input: `cat /dev/null`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1124",
					Message: "Use `: > file` instead of `cat /dev/null > file` to truncate. The `:` builtin avoids spawning cat.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1124")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
