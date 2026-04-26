// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1301(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid pipestatus usage",
			input:    `echo $pipestatus`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid PIPESTATUS usage",
			input: `echo $PIPESTATUS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1301",
					Message: "Avoid `$PIPESTATUS` in Zsh — use `$pipestatus` (lowercase) instead. The uppercase form is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1301")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
