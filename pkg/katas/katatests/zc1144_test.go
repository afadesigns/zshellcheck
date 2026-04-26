// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1144(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid trap with name",
			input:    `trap cleanup EXIT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid trap with 0",
			input:    `trap cleanup 0`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid trap with number",
			input: `trap cleanup 15`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1144",
					Message: "Use signal names (`SIGTERM`, `SIGINT`, `EXIT`) instead of numbers in `trap`. Signal numbers vary across platforms.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1144")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
