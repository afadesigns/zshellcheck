// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1306(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid CURRENT usage",
			input:    `echo $CURRENT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid COMP_CWORD usage",
			input: `echo $COMP_CWORD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1306",
					Message: "Avoid `$COMP_CWORD` in Zsh — use `$CURRENT` instead. `COMP_CWORD` is Bash completion-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1306")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
