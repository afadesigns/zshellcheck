// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1348(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -type",
			input:    `find . -name "*.txt"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -type f",
			input: `find . -type f`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1348",
					Message: "Use Zsh glob type qualifiers (`*(/)`, `*(.)`, `*(@)`, `*(=)`, `*(p)`, `*(*)`, `*(%)`) instead of `find -type`. No external process required.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -type d",
			input: `find / -type d`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1348",
					Message: "Use Zsh glob type qualifiers (`*(/)`, `*(.)`, `*(@)`, `*(=)`, `*(p)`, `*(*)`, `*(%)`) instead of `find -type`. No external process required.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1348")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
