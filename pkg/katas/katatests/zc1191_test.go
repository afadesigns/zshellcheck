// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1191(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid other command",
			input:    `reset`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid bare clear",
			input: `clear`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1191",
					Message: "Use `print -n '\\e[2J\\e[H'` instead of `clear`. ANSI escape sequences avoid spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1191")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
