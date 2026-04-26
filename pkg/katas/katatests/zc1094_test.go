// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1094(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sed with -i flag",
			input:    `sed -i 's/foo/bar/' file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sed with file argument",
			input:    `sed 's/foo/bar/' file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sed with -e flag",
			input:    `sed -e 's/foo/bar/' -e 's/baz/qux/'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid simple sed substitution",
			input: `sed 's/foo/bar/'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1094",
					Message: "Use `${var//pattern/replacement}` instead of piping through `sed` for simple substitutions. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid sed global substitution",
			input: `sed 's/foo/bar/g'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1094",
					Message: "Use `${var//pattern/replacement}` instead of piping through `sed` for simple substitutions. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1094")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
