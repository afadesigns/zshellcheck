// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1133(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid kill SIGTERM",
			input:    `kill 1234`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid kill -15",
			input:    `kill -15 1234`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid kill -9",
			input: `kill -9 1234`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1133",
					Message: "Avoid `kill -9` as a first resort. Use `kill` (SIGTERM) first to allow graceful shutdown, then escalate to `kill -9` if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid kill -KILL",
			input: `kill -KILL 1234`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1133",
					Message: "Avoid `kill -9` as a first resort. Use `kill` (SIGTERM) first to allow graceful shutdown, then escalate to `kill -9` if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1133")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
