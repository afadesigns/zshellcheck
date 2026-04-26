// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1871(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt IGNORE_BRACES` (explicit default)",
			input:    `unsetopt IGNORE_BRACES`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt IGNORE_BRACES`",
			input: `setopt IGNORE_BRACES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1871",
					Message: "`setopt IGNORE_BRACES` disables brace expansion — `for i in {1..10}` loops once over the literal token, `cp app.{conf,bak}` fails ENOENT. Keep the option off; quote the specific argument if you need a literal brace string.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_IGNORE_BRACES`",
			input: `unsetopt NO_IGNORE_BRACES`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1871",
					Message: "`unsetopt NO_IGNORE_BRACES` disables brace expansion — `for i in {1..10}` loops once over the literal token, `cp app.{conf,bak}` fails ENOENT. Keep the option off; quote the specific argument if you need a literal brace string.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1871")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
