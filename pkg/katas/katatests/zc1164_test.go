// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1164(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sed substitution",
			input:    `sed -n 's/foo/bar/p'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sed with file",
			input:    `sed -n '3p' file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sed -n Np in pipeline",
			input: `sed -n '5p'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1164",
					Message: "Use Zsh array subscript `${lines[N]}` instead of `sed -n 'Np'`. Split input with `${(f)...}` then index directly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1164")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
