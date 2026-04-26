// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1302(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid man usage",
			input:    `man zshbuiltins`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid help usage",
			input: `help cd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1302",
					Message: "Avoid `help` in Zsh — it is a Bash builtin. Use `run-help` or `man zshbuiltins` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1302")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
