package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1743(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `npm audit fix` (no --force)",
			input:    `npm audit fix`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `npm audit` (no fix)",
			input:    `npm audit`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `npm run build --force` (not audit fix)",
			input:    `npm run build --force`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `npm audit fix --force`",
			input: `npm audit fix --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1743",
					Message: "`npm audit ... --force` accepts every major-version bump an advisory triggers — silent breaking changes. Drop `--force` and triage advisories one by one.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pnpm audit --fix --force`",
			input: `pnpm audit --fix --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1743",
					Message: "`pnpm audit ... --force` accepts every major-version bump an advisory triggers — silent breaking changes. Drop `--force` and triage advisories one by one.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1743")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
