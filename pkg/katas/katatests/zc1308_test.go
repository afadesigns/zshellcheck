// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1308(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid BUFFER usage",
			input:    `echo $BUFFER`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid COMP_LINE usage",
			input: `echo $COMP_LINE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1308",
					Message: "Avoid `$COMP_LINE` in Zsh — use `$BUFFER` instead. `COMP_LINE` is Bash completion-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1308")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
