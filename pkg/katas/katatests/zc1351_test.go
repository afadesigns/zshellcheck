// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1351(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — expr arithmetic",
			input:    `expr 1 + 2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — expr match",
			input: `expr match "$s" '^foo'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1351",
					Message: "Use `[[ $str =~ pattern ]]` with `$match` / `$MATCH` arrays instead of `expr match`/`expr index`. Regex evaluation stays in the shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — expr index",
			input: `expr index "$s" aeiou`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1351",
					Message: "Use `[[ $str =~ pattern ]]` with `$match` / `$MATCH` arrays instead of `expr match`/`expr index`. Regex evaluation stays in the shell.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1351")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
