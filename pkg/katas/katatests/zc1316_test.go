// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1316(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid funcfiletrace usage",
			input:    `echo $funcfiletrace`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid caller usage",
			input: `caller 0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1316",
					Message: "Avoid `caller` in Zsh — it is a Bash builtin. Use `$funcfiletrace` and `$funcstack` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1316")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
