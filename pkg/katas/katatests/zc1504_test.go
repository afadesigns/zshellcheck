// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1504(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — git push origin main",
			input:    `git push origin main`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — git push --all origin",
			input:    `git push --all origin`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — git push --mirror origin",
			input: `git push --mirror origin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1504",
					Message: "`git push --mirror` overwrites every remote ref and deletes ones missing locally. Use an explicit refspec or `--all` for everyday pushes.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — git push origin --mirror",
			input: `git push origin --mirror`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1504",
					Message: "`git push --mirror` overwrites every remote ref and deletes ones missing locally. Use an explicit refspec or `--all` for everyday pushes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1504")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
