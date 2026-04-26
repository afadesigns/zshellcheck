// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1307(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid dirstack usage",
			input:    `echo $dirstack`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid DIRSTACK usage",
			input: `echo $DIRSTACK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1307",
					Message: "Avoid `$DIRSTACK` in Zsh — use `$dirstack` (lowercase) instead. The uppercase form is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1307")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
