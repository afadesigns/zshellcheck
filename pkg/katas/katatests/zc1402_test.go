// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1402(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — date +fmt",
			input:    `date +%Y-%m-%d`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — date -d",
			input: `date -d @1700000000 +%Y-%m-%d`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1402",
					Message: "Use Zsh `strftime` (from `zsh/datetime`) instead of `date -d @N -- +fmt`. The `-d`/`@` form is GNU-specific; `strftime` is portable Zsh.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1402")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
