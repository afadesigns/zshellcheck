package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1509(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — trap 'cleanup' EXIT",
			input:    `trap 'cleanup' EXIT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — trap 'cleanup' TERM",
			input:    `trap 'cleanup' TERM`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — trap '' TERM",
			input: `trap '' TERM`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1509",
					Message: "`trap '' TERM` silences a fatal signal — cleanup handlers never run. Keep at least a cleanup trap on EXIT.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — trap - SIGINT",
			input: `trap - SIGINT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1509",
					Message: "`trap - SIGINT` silences a fatal signal — cleanup handlers never run. Keep at least a cleanup trap on EXIT.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1509")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
