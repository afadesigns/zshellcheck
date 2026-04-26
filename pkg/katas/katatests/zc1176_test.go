// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1176(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid zparseopts",
			input:    `zparseopts -D -E -- v=verbose h=help`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid getopt",
			input: `getopt -o vh -l verbose,help`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1176",
					Message: "Use Zsh `zparseopts` instead of `getopt`. `zparseopts` supports long options, arrays, and is the native Zsh approach.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid getopts",
			input: `getopts "vh" opt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1176",
					Message: "Use Zsh `zparseopts` instead of `getopts`. `zparseopts` supports long options, arrays, and is the native Zsh approach.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1176")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
