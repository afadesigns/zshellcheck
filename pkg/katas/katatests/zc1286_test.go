// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1286(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep without -v",
			input:    `grep pattern file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid Zsh array filter",
			input:    `echo ${array:#pattern}`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep -v for filtering",
			input: `grep -v pattern file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1286",
					Message: "Use Zsh `${array:#pattern}` for filtering instead of `grep -v`. Parameter expansion avoids a subprocess.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1286")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
