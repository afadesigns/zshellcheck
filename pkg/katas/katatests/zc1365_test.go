// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1365(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — stat without format flag",
			input:    `stat file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — stat -c %s",
			input: `stat -c %s file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1365",
					Message: "Use Zsh `zmodload zsh/stat; zstat -H meta file` for file metadata instead of `stat -c '%...'`. The associative array `meta` exposes every stat field.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — stat -c %Y (mtime)",
			input: `stat -c %Y file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1365",
					Message: "Use Zsh `zmodload zsh/stat; zstat -H meta file` for file metadata instead of `stat -c '%...'`. The associative array `meta` exposes every stat field.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1365")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
