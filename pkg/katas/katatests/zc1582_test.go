// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1582(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — bash script.sh",
			input:    `bash script.sh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — bash -x script.sh",
			input: `bash -x script.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1582",
					Message: "`bash -x` traces every expanded command — CI logs leak secrets verbatim. Scope with `set -x; …; set +x`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sh -x script.sh",
			input: `sh -x script.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1582",
					Message: "`sh -x` traces every expanded command — CI logs leak secrets verbatim. Scope with `set -x; …; set +x`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — zsh -xv script.zsh",
			input: `zsh -xv script.zsh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1582",
					Message: "`zsh -xv` traces every expanded command — CI logs leak secrets verbatim. Scope with `set -x; …; set +x`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1582")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
