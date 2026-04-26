// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1293(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid [[ ]] usage",
			input:    `[[ -f file.txt ]]`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid test command",
			input: `test -f file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1293",
					Message: "Use `[[ ]]` instead of the `test` command in Zsh. `[[ ]]` is more powerful and does not require variable quoting.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid test with -z flag",
			input: `test -z "$var"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1293",
					Message: "Use `[[ ]]` instead of the `test` command in Zsh. `[[ ]]` is more powerful and does not require variable quoting.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1293")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
