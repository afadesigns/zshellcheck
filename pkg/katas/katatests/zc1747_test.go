package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1747(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `npm install --registry=https://registry.npmjs.org/`",
			input:    `npm install --registry=https://registry.npmjs.org/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `yarn config set registry https://registry.npmjs.org/`",
			input:    `yarn config set registry https://registry.npmjs.org/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `npm install --registry=http://internal/`",
			input: `npm install --registry=http://internal/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1747",
					Message: "`npm --registry=http://internal/` uses plaintext HTTP for the package registry — any MITM swaps tarballs and runs install-time `postinstall` code. Use `https://`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pnpm install --registry http://internal/`",
			input: `pnpm install --registry http://internal/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1747",
					Message: "`pnpm --registry http://internal/` uses plaintext HTTP for the package registry — any MITM swaps tarballs and runs install-time `postinstall` code. Use `https://`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `yarn config set registry http://internal/`",
			input: `yarn config set registry http://internal/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1747",
					Message: "`yarn config set registry http://internal/` uses plaintext HTTP for the package registry — any MITM swaps tarballs and runs install-time `postinstall` code. Use `https://`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1747")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
