package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1553(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — tr -d '[:space:]'",
			input:    `tr -d '[:space:]'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — tr -s ' '",
			input:    `tr -s ' '`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tr '[:lower:]' '[:upper:]'",
			input: `tr '[:lower:]' '[:upper:]'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1553",
					Message: "`tr` for case conversion — use Zsh `${(U)var}` / `${(L)var}` to avoid the fork/exec and portability hazard.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tr a-z A-Z",
			input: `tr a-z A-Z`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1553",
					Message: "`tr` for case conversion — use Zsh `${(U)var}` / `${(L)var}` to avoid the fork/exec and portability hazard.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — tr '[:upper:]' '[:lower:]'",
			input: `tr '[:upper:]' '[:lower:]'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1553",
					Message: "`tr` for case conversion — use Zsh `${(U)var}` / `${(L)var}` to avoid the fork/exec and portability hazard.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1553")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
