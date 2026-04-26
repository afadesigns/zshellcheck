// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1817(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `git push origin main`",
			input:    `git push origin main`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `git push -u origin feature-x`",
			input:    `git push -u origin feature-x`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `git push --delete origin mybranch`",
			input: `git push --delete origin mybranch`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1817",
					Message: "`git push --delete` deletes the remote branch — open PRs are orphaned, CI targets disappear, and the last commit SHA can only come back from someone else's clone. Let the hosting platform auto-delete after merge instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `git push origin :mybranch`",
			input: `git push origin :mybranch`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1817",
					Message: "`git push origin :mybranch` deletes the remote branch — open PRs are orphaned, CI targets disappear, and the last commit SHA can only come back from someone else's clone. Let the hosting platform auto-delete after merge instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1817")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
