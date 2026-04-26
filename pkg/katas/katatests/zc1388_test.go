// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1388(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $mailpath (Zsh)",
			input:    `echo $mailpath`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $MAILPATH",
			input: `echo $MAILPATH`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1388",
					Message: "Use Zsh lowercase `$mailpath` (array) instead of Bash uppercase `$MAILPATH` (colon-separated string).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1388")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
