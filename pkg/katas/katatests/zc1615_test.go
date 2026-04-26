// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1615(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — plain date",
			input:    `date "+%Y-%m-%d"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — date +%s only (ZC1119 handles)",
			input:    `date "+%s"`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — date "+%s.%N"`,
			input: `date "+%s.%N"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1615",
					Message: "`date \"+%s.%N\"` forks for sub-second time. Use Zsh `$EPOCHREALTIME` / `$epochtime` from `zmodload zsh/datetime`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — date "+%s%N"`,
			input: `date "+%s%N"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1615",
					Message: "`date \"+%s%N\"` forks for sub-second time. Use Zsh `$EPOCHREALTIME` / `$epochtime` from `zmodload zsh/datetime`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1615")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
