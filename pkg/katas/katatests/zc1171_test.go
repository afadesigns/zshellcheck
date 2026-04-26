// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1171(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid print",
			input:    `print "hello\nworld"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid echo without -e",
			input:    `echo "hello"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid echo -e",
			input: `echo -e "hello\nworld"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1171",
					Message: "Use `print` instead of `echo -e`. Zsh `print` natively interprets escape sequences and is more portable than `echo -e`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1171")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
