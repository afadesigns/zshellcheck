// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1319(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid argument count",
			input:    `echo $MYVAR`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_ARGC usage",
			input: `echo $BASH_ARGC`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1319",
					Message: "Avoid `$BASH_ARGC` in Zsh — use `$#` for argument count. `BASH_ARGC` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1319")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
