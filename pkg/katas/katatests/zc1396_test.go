package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1396(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unset -v var",
			input:    `unset -v var`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — unset -n ref",
			input: `unset -n ref`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1396",
					Message: "`unset -n` is a Bash nameref operation. Zsh does not honor it; use `unset -v NAME` (variable) or `unset -f NAME` (function) explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1396")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
