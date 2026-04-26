// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1375(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — tty (print tty name)",
			input:    `tty`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tty -s",
			input: `tty -s`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1375",
					Message: "Use `[[ -t 0 ]]` (stdin), `[[ -t 1 ]]` (stdout), or `[[ -t 2 ]]` (stderr) instead of `tty -s`. In-shell file-descriptor test, no external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1375")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
