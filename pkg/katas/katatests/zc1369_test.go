package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1369(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — od -x (hex, different use)",
			input:    `od -x file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — od -c",
			input: `od -c file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1369",
					Message: "Use Zsh `${(V)var}` to see non-printable characters in a variable — renders control chars as `\\n`, `\\t`, etc., without spawning `od`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1369")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
