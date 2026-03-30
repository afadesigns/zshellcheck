package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1061(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "seq usage",
			input: `seq 1 10`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1061",
					Message: "Prefer `{start..end}` range expansion over `seq`. It is built-in and faster.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "no seq",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1061")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
