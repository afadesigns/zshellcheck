// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1152(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -E",
			input:    `grep -E "pattern" file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep -P",
			input: `grep -P "\d+" file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1152",
					Message: "Avoid `grep -P` — it's unavailable on macOS. Use `zmodload zsh/pcre` with `pcre_compile`/`pcre_match` or `grep -E` for portable regex matching.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1152")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
