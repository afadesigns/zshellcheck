package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1963(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `npx typescript@5.4.2 --init`",
			input:    `npx typescript@5.4.2 --init`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pnpm dlx @vercel/ncc@0.38.1 build`",
			input:    `pnpm dlx @vercel/ncc@0.38.1 build`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `npx create-react-app demo`",
			input: `npx create-react-app demo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1963",
					Message: "`npx create-react-app` pulls the `latest` tag every run — a squatted or compromised package lands attacker code. Pin the version (`create-react-app@X.Y.Z`) or use a regular `npm install` + `./node_modules/.bin/`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pnpm dlx prettier`",
			input: `pnpm dlx prettier`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1963",
					Message: "`pnpm dlx prettier` pulls the `latest` tag every run — a squatted or compromised package lands attacker code. Pin the version (`prettier@X.Y.Z`) or use a regular `npm install` + `./node_modules/.bin/`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1963")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
