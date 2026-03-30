package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1177(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid id -g",
			input:    `id -g`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid id -u",
			input: `id -u`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1177",
					Message: "Use `$UID` or `$EUID` instead of `id -u`. Zsh provides these as built-in variables.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1177")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
