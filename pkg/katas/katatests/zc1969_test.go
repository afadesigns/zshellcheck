// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1969(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `zsh -c ':'` (no -f/-d flag)",
			input:    `zsh -c ':'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `zsh $SCRIPT`",
			input:    `zsh $SCRIPT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `zsh -f $SCRIPT`",
			input: `zsh -f $SCRIPT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1969",
					Message: "`zsh -f` skips `/etc/zsh*` and `~/.zsh*` startup files — corp proxy/audit/`PATH` hardening silently dropped. For a pristine shell use `env -i zsh` with an explicit allow-list.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `zsh -d $SCRIPT`",
			input: `zsh -d $SCRIPT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1969",
					Message: "`zsh -d` skips `/etc/zsh*` and `~/.zsh*` startup files — corp proxy/audit/`PATH` hardening silently dropped. For a pristine shell use `env -i zsh` with an explicit allow-list.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1969")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
