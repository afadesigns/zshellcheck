// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1589(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — trap cleanup EXIT",
			input:    `trap cleanup EXIT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — trap with safe dump",
			input:    `trap 'echo failed' ERR`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — trap 'set -x' ERR",
			input: `trap 'set -x' ERR`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1589",
					Message: "`trap 'set -x' ERR` enables shell trace from a trap — expansions leak env vars (tokens, passwords) to stderr / CI logs. Use a scoped `set -x ... set +x`, not a trap.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — trap 'set -o xtrace' EXIT",
			input: `trap 'set -o xtrace' EXIT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1589",
					Message: "`trap 'set -x' EXIT` enables shell trace from a trap — expansions leak env vars (tokens, passwords) to stderr / CI logs. Use a scoped `set -x ... set +x`, not a trap.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — trap 'set -x' RETURN",
			input: `trap 'set -x' RETURN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1589",
					Message: "`trap 'set -x' RETURN` enables shell trace from a trap — expansions leak env vars (tokens, passwords) to stderr / CI logs. Use a scoped `set -x ... set +x`, not a trap.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1589")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
