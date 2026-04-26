// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1335(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid other command",
			input:    `cat file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "tac usage",
			input: `tac file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1335",
					Message: "Consider Zsh `${(Oa)array}` for reversing array data instead of piping to `tac`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1335")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
