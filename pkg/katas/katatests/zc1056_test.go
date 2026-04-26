// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1056(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "normal command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "command substitution as statement",
			input:    `$(echo hello)`,
			expected: []katas.Violation{},
		},
		{
			name:     "arithmetic command",
			input:    `((x + 1))`,
			expected: []katas.Violation{},
		},
		{
			name:     "variable expansion command",
			input:    `$var`,
			expected: []katas.Violation{},
		},
		{
			name:     "ls command no violation",
			input:    `ls -la`,
			expected: []katas.Violation{},
		},
		{
			name:     "assignment statement",
			input:    `x=5`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1056")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
