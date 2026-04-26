// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1386(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated echo",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo $FIGNORE",
			input: `echo $FIGNORE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1386",
					Message: "`$FIGNORE` is Bash-only. In Zsh use `zstyle ':completion:*' ignored-patterns '*.o *.pyc'` for completion filtering.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1386")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
