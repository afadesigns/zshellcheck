// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1212(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid git add file",
			input:    `git add main.go`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid git add dot",
			input: `git add .`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1212",
					Message: "Avoid `git add .` or `git add -A` — they stage everything including unintended files. Use explicit paths or `git add -p` for selective staging.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid git add -A",
			input: `git add -A`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1212",
					Message: "Avoid `git add .` or `git add -A` — they stage everything including unintended files. Use explicit paths or `git add -p` for selective staging.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1212")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
