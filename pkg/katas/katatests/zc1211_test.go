// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1211(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid git stash push -m",
			input:    `git stash push -m "wip"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid git stash pop",
			input:    `git stash pop`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid bare git stash",
			input: `git stash`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1211",
					Message: "Use `git stash push -m 'description'` instead of bare `git stash`. Named stashes are easier to identify and manage.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1211")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
