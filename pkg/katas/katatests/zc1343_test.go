// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1343(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without age predicate",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -mtime +7",
			input: `find . -mtime +7`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1343",
					Message: "Use Zsh glob qualifiers (`*(m±N)`, `*(M±N)`, `*(a±N)`, `*(c±N)`) instead of `find -mtime`/`-mmin`/`-atime`/`-amin`/`-ctime`/`-cmin` for age predicates.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -mmin -60",
			input: `find . -mmin -60`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1343",
					Message: "Use Zsh glob qualifiers (`*(m±N)`, `*(M±N)`, `*(a±N)`, `*(c±N)`) instead of `find -mtime`/`-mmin`/`-atime`/`-amin`/`-ctime`/`-cmin` for age predicates.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -ctime 0",
			input: `find . -ctime 0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1343",
					Message: "Use Zsh glob qualifiers (`*(m±N)`, `*(M±N)`, `*(a±N)`, `*(c±N)`) instead of `find -mtime`/`-mmin`/`-atime`/`-amin`/`-ctime`/`-cmin` for age predicates.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1343")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
