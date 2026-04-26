// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1368(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sh without -c",
			input:    `sh script.sh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sh -c",
			input: `sh -c 'echo hi'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1368",
					Message: "Avoid `sh -c`/`bash -c` inside a Zsh script — inline the code as a function to keep access to arrays, associative arrays, and Zsh features. Use `zsh -c` only when a fresh shell is truly required.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — bash -c",
			input: `bash -c 'echo hi'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1368",
					Message: "Avoid `sh -c`/`bash -c` inside a Zsh script — inline the code as a function to keep access to arrays, associative arrays, and Zsh features. Use `zsh -c` only when a fresh shell is truly required.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1368")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
