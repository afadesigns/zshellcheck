// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1303(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid zmodload usage",
			input:    `zmodload zsh/stat`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid enable without -f",
			input:    `enable mybuiltin`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid enable -f usage",
			input: `enable -f /path/to/builtin mybuiltin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1303",
					Message: "Avoid `enable -f` in Zsh — use `zmodload` to load modules. `enable -f` is Bash-specific.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1303")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
