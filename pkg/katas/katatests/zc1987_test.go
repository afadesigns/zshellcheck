// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1987(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt BRACE_CCL` (default)",
			input:    `unsetopt BRACE_CCL`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NO_BRACE_CCL`",
			input:    `setopt NO_BRACE_CCL`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt BRACE_CCL`",
			input: `setopt BRACE_CCL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1987",
					Message: "`setopt BRACE_CCL` promotes single-character braces to csh-style classes — `{a-z}` now expands to every letter, `{ABC}` to `A B C`, breaking regex/hex/CI-name literals. Use `{a..z}` when a real range is wanted.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_BRACE_CCL`",
			input: `unsetopt NO_BRACE_CCL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1987",
					Message: "`unsetopt NO_BRACE_CCL` promotes single-character braces to csh-style classes — `{a-z}` now expands to every letter, `{ABC}` to `A B C`, breaking regex/hex/CI-name literals. Use `{a..z}` when a real range is wanted.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1987")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
