package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1163(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -m 1",
			input:    `grep -m 1 pattern file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep | head -1",
			input: `grep pattern file | head -1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1163",
					Message: "Use `grep -m 1` instead of `grep | head -1`. The `-m` flag stops after the first match without a pipeline.",
					Line:    1,
					Column:  19,
				},
			},
		},
		{
			name:     "non-pipe operator",
			input:    `echo hello && echo world`,
			expected: []katas.Violation{},
		},
		{
			name:     "pipe but left not grep",
			input:    `cat file | head -1`,
			expected: []katas.Violation{},
		},
		{
			name:     "grep piped to non-head",
			input:    `grep pattern file | sort`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1163")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
