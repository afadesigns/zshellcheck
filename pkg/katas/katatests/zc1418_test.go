// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1418(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ulimit -t",
			input:    `ulimit -t 10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ulimit -H",
			input: `ulimit -H -t 60`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1418",
					Message: "Use Zsh `limit -h` (hard) / `limit -s` (soft) instead of `ulimit -H`/`-S`. Zsh's `limit` builtin is more human-readable.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1418")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
