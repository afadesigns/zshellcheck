// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1282(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid Zsh :r modifier",
			input:    `echo ${file:r}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sed with different pattern",
			input:    `sed s/foo/bar/g`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sed to strip extension",
			input: `sed 's/\.[^.]*$//'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1282",
					Message: "Use Zsh parameter expansion `${var:r}` to remove the file extension instead of `sed`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1282")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
