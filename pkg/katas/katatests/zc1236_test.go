// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1236(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid git reset --soft",
			input:    `git reset --soft HEAD~1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid git reset --hard",
			input: `git reset --hard HEAD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1236",
					Message: "Avoid `git reset --hard` — it permanently discards uncommitted changes. Use `git stash` first, or `git reset --soft` to keep changes staged.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1236")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
