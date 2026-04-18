package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1640(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — Zsh (P) flag",
			input:    `echo "${(P)var}"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — plain expansion",
			input:    `echo "${var}"`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — echo "${!var}"`,
			input: `echo "${!var}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1640",
					Message: "`${!var}` Bash indirect — prefer Zsh `${(P)var}` for the same semantics with flag composability.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — print "${!array[@]}"`,
			input: `print -r -- "${!array[@]}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1640",
					Message: "`${!var}` Bash indirect — prefer Zsh `${(P)var}` for the same semantics with flag composability.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1640")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
