// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1530(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pkill sshd (process name)",
			input:    `pkill sshd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — pkill -U 1000 java",
			input:    `pkill -U 1000 java`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pkill -f server",
			input: `pkill -f server`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1530",
					Message: "`pkill -f` matches the full command line — easy to over-kill. Drop `-f`, scope with `-U/-G/-P`, or anchor the pattern with ^/$.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1530")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
