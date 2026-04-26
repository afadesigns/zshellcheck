// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1322(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid non-COPROC variable",
			input:    `echo $MY_PROC`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid COPROC usage",
			input: `echo $COPROC`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1322",
					Message: "Avoid `$COPROC` in Zsh — Zsh coprocesses use `read -p`/`print -p` for I/O. `COPROC` is Bash-specific.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1322")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
