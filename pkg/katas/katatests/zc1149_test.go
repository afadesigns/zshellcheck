// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1149(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid echo normal message",
			input:    `echo "Processing..."`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid echo error to stdout",
			input: `echo "Error: file not found"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1149",
					Message: "Error messages should go to stderr. Use `print -u2` or append `>&2` to separate error output from normal stdout.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1149")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
