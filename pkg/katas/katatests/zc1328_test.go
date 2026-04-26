// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1328(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid HISTSIZE usage",
			input:    `echo $HISTSIZE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid HISTCONTROL usage",
			input: `echo $HISTCONTROL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1328",
					Message: "Avoid `$HISTCONTROL` in Zsh — use `setopt HIST_IGNORE_DUPS` and related options instead.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1328")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
