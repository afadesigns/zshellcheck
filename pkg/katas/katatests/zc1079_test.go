// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1079(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid quoted comparison",
			input:    `[[ $var == "$other" ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid literal comparison",
			input:    `[[ $var == "foo" ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid pattern comparison (literal)",
			input:    `[[ $var == foo* ]]`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid unquoted variable == ",
			input: `[[ $var == $other ]]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1079",
					Message: "Unquoted RHS matches as pattern. Quote to force string comparison: `\"$var\"`.",
					Line:    1,
					Column:  12,
				},
			},
		},
		{
			name:  "invalid unquoted variable !=",
			input: `[[ $var != $other ]]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1079",
					Message: "Unquoted RHS matches as pattern. Quote to force string comparison: `\"$var\"`.",
					Line:    1,
					Column:  12,
				},
			},
		},
		{
			name:  "invalid array access",
			input: `[[ $var = ${arr[1]} ]]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1079",
					Message: "Unquoted RHS matches as pattern. Quote to force string comparison: `\"$var\"`.",
					Line:    1,
					Column:  11,
				},
			},
		},
		{
			name:     "non-equality operator",
			input:    `[[ $a -lt $b ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "quoted string RHS",
			input:    `[[ $a == "hello" ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "integer comparison",
			input:    `[[ $a -eq 5 ]]`,
			expected: []katas.Violation{},
		},
		{
			name:     "no infix elements",
			input:    `[[ -f /tmp/foo ]]`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1079")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
