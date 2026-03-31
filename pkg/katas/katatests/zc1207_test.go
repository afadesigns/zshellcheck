package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1207(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid chpasswd",
			input:    `chpasswd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid passwd",
			input: `passwd user`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1207",
					Message: "Avoid `passwd` in scripts — it prompts interactively. Use `chpasswd` or `usermod --password` for non-interactive password changes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1207")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
