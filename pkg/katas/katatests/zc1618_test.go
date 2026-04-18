package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1618(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — git commit -m (no skip)",
			input:    `git commit -m "msg"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — git push --dry-run",
			input:    `git push --dry-run origin main`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — git commit --no-verify`,
			input: `git commit --no-verify -m "msg"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1618",
					Message: "`git commit --no-verify` skips pre-commit / commit-msg hooks — lint, test, and secret-scan checks do not run. Reserve for emergencies; scripts should pass the hooks.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — git push --no-verify`,
			input: `git push --no-verify origin main`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1618",
					Message: "`git push --no-verify` skips pre-push / commit-msg hooks — lint, test, and secret-scan checks do not run. Reserve for emergencies; scripts should pass the hooks.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — git commit -n -m`,
			input: `git commit -n -m "msg"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1618",
					Message: "`git commit -n` skips pre-commit / commit-msg hooks — lint, test, and secret-scan checks do not run. Reserve for emergencies; scripts should pass the hooks.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1618")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
