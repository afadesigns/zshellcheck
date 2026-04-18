package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1696(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pnpm install --frozen-lockfile",
			input:    `pnpm install --frozen-lockfile`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — npm ci",
			input:    `npm ci`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pnpm install --no-frozen-lockfile",
			input: `pnpm install --no-frozen-lockfile`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1696",
					Message: "`--no-frozen-lockfile` allows the lockfile to drift — the CI artifact no longer matches the reviewed dependency graph. Use `--frozen-lockfile` / `--immutable` in CI.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — yarn install --no-immutable",
			input: `yarn install --no-immutable`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1696",
					Message: "`--no-immutable` allows the lockfile to drift — the CI artifact no longer matches the reviewed dependency graph. Use `--frozen-lockfile` / `--immutable` in CI.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1696")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
