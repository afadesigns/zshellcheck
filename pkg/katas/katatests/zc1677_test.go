// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1677(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — trap cleanup EXIT",
			input:    `trap 'cleanup' EXIT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — trap set -x on ERR (different signal)",
			input:    `trap 'set -x' ERR`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — trap 'set -x' DEBUG",
			input: `trap 'set -x' DEBUG`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1677",
					Message: "`trap 'set -x' DEBUG` keeps xtrace on after the first command — every subsequent argv (passwords, bearer tokens) lands in the log. Trace a narrow block with `set -x … set +x` or use `typeset -ft FUNC` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — trap "set -o xtrace" DEBUG`,
			input: `trap "set -o xtrace" DEBUG`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1677",
					Message: "`trap 'set -x' DEBUG` keeps xtrace on after the first command — every subsequent argv (passwords, bearer tokens) lands in the log. Trace a narrow block with `set -x … set +x` or use `typeset -ft FUNC` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1677")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
