// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1379(t *testing.T) {
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
			name:  "invalid — echo $PROMPT_COMMAND",
			input: `echo $PROMPT_COMMAND`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1379",
					Message: "`PROMPT_COMMAND` is Bash-only. In Zsh define a `precmd` function or use `autoload -Uz add-zsh-hook; add-zsh-hook precmd my_hook`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1379")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
