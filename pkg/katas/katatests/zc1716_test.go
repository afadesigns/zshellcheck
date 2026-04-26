// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1716(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — Zsh-idiomatic `$CPUTYPE`",
			input:    `print -r -- $CPUTYPE`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `uname -r` (kernel release, no Zsh equivalent)",
			input:    `uname -r`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `uname -m`",
			input: `uname -m`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1716",
					Message: "Use Zsh `$CPUTYPE` / `$MACHTYPE` instead of `uname -m` — parameter expansion avoids forking an external for an answer Zsh already cached at startup.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `uname -p`",
			input: `uname -p`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1716",
					Message: "Use Zsh `$CPUTYPE` / `$MACHTYPE` instead of `uname -p` — parameter expansion avoids forking an external for an answer Zsh already cached at startup.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1716")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
