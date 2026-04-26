// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1477(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — printf '%s\\n' \"$x\"",
			input:    `printf '%s\n' "$x"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — printf \"%s\\n\" \"$x\"",
			input:    `printf "%s\n" "$x"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — printf 'hello world'",
			input:    `printf 'hello world'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — printf \"$x\"",
			input: `printf "$x"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1477",
					Message: "`printf` format string contains a variable — `%` inside `$var` is reparsed as a format specifier. Use `printf '%s' \"$var\"` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — printf \"$(cmd)\"",
			input: `printf "$(cmd)"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1477",
					Message: "`printf` format string contains a variable — `%` inside `$var` is reparsed as a format specifier. Use `printf '%s' \"$var\"` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1477")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
