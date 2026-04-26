// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1997(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt HIST_NO_FUNCTIONS` (default)",
			input:    `unsetopt HIST_NO_FUNCTIONS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_HIST_NO_FUNCTIONS`",
			input:    `setopt NO_HIST_NO_FUNCTIONS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt HIST_NO_FUNCTIONS`",
			input: `setopt HIST_NO_FUNCTIONS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1997",
					Message: "`setopt HIST_NO_FUNCTIONS` drops function-definition commands from `$HISTFILE` — forensic trail loses the definition while the call that used it still shows. Scope hiding via `zshaddhistory` hook instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_HIST_NO_FUNCTIONS`",
			input: `unsetopt NO_HIST_NO_FUNCTIONS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1997",
					Message: "`unsetopt NO_HIST_NO_FUNCTIONS` drops function-definition commands from `$HISTFILE` — forensic trail loses the definition while the call that used it still shows. Scope hiding via `zshaddhistory` hook instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1997")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
