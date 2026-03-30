package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1130(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid true",
			input: `true`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1130",
					Message: "Use `:` instead of `true`. `:` is always a shell builtin, while `true` may be an external command.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid colon",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1130")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
