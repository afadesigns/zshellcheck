package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1452(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — npm install local",
			input:    `npm install react`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — npm install -g",
			input: `npm install -g typescript`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1452",
					Message: "`npm install -g` installs system-wide. Prefer project-local install or `npx`/`pnpm dlx` for one-off tools.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1452")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
