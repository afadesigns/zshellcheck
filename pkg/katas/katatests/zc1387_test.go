// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1387(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $options (Zsh)",
			input:    `echo $options`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $SHELLOPTS",
			input: `echo $SHELLOPTS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1387",
					Message: "`$SHELLOPTS` is Bash-only. In Zsh inspect `$options` (assoc array, keys are option names) via `print -l ${(kv)options}`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1387")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
