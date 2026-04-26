// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1877(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt SHORT_LOOPS` (explicit default)",
			input:    `setopt SHORT_LOOPS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt SHORT_LOOPS`",
			input: `unsetopt SHORT_LOOPS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1877",
					Message: "`unsetopt SHORT_LOOPS` disables short-form loops — `for f in *.log; print $f` raises a parse error. Keep the option on; scope inside a function with `LOCAL_OPTIONS` if POSIX-strict parsing is really needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_SHORT_LOOPS`",
			input: `setopt NO_SHORT_LOOPS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1877",
					Message: "`setopt NO_SHORT_LOOPS` disables short-form loops — `for f in *.log; print $f` raises a parse error. Keep the option on; scope inside a function with `LOCAL_OPTIONS` if POSIX-strict parsing is really needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1877")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
