// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1398(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $PS1",
			input:    `echo $PS1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $PROMPT_DIRTRIM",
			input: `echo $PROMPT_DIRTRIM`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1398",
					Message: "`$PROMPT_DIRTRIM` is Bash-only. Use the Zsh prompt escape `%N~` (N = number of path components to keep) for directory truncation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1398")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
