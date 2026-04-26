// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1294(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid bindkey usage",
			input:    `bindkey -e`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid bind usage",
			input: `bind -x '"\C-r": history-search'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1294",
					Message: "Use `bindkey` instead of `bind` in Zsh. `bind` is a Bash builtin; Zsh uses `bindkey` for ZLE key bindings.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1294")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
