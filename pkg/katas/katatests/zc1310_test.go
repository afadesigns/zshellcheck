// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1310(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid variable",
			input:    `echo $MY_STRING`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid BASH_EXECUTION_STRING usage",
			input: `echo $BASH_EXECUTION_STRING`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1310",
					Message: "Avoid `$BASH_EXECUTION_STRING` in Zsh — it is undefined. Access command arguments directly instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1310")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
