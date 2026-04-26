// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1720(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — Zsh-idiomatic `$COLUMNS`",
			input:    `print -r -- $COLUMNS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `tput setaf 1` (color, no $COLUMNS equivalent)",
			input:    `tput setaf 1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tput cols`",
			input: `tput cols`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1720",
					Message: "Use `$COLUMNS` instead of `tput cols` — Zsh keeps the terminal size in parameters, no fork needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `tput lines`",
			input: `tput lines`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1720",
					Message: "Use `$LINES` instead of `tput lines` — Zsh keeps the terminal size in parameters, no fork needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1720")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
