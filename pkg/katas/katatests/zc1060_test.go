package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1060(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "no pipe",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "pipe without ps",
			input:    `cat file | grep pattern`,
			expected: []katas.Violation{},
		},
		{
			name:  "ps piped to grep",
			input: `ps aux | grep myprocess`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1060",
					Message: "`ps | grep pattern` matches the grep process itself. Use `grep [p]attern` to exclude the grep process.",
					Line:    1,
					Column:  8,
				},
			},
		},
		{
			name:     "ps piped to non-grep command",
			input:    `ps aux | sort`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1060")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
