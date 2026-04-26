// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1923(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt PRINT_EXIT_VALUE` (explicit default)",
			input:    `unsetopt PRINT_EXIT_VALUE`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt EXTENDED_HISTORY` (unrelated)",
			input:    `setopt EXTENDED_HISTORY`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt PRINT_EXIT_VALUE`",
			input: `setopt PRINT_EXIT_VALUE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1923",
					Message: "`setopt PRINT_EXIT_VALUE` prints `zsh: exit N` on stderr for every non-zero exit — silent grep/test/curl probes suddenly leak status, and tools parsing stderr see interleaved shell chatter. Remove; use `|| printf …` per call.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_PRINT_EXIT_VALUE`",
			input: `unsetopt NO_PRINT_EXIT_VALUE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1923",
					Message: "`unsetopt NO_PRINT_EXIT_VALUE` prints `zsh: exit N` on stderr for every non-zero exit — silent grep/test/curl probes suddenly leak status, and tools parsing stderr see interleaved shell chatter. Remove; use `|| printf …` per call.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1923")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
