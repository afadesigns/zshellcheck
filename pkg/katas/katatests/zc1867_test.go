// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1867(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt GLOB` (explicit default)",
			input:    `setopt GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt GLOB`",
			input: `unsetopt GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1867",
					Message: "`unsetopt GLOB` disables glob expansion — `rm *.log` chases the literal `*.log`, `for f in *.txt` loops once. Quote specific args or scope with `LOCAL_OPTIONS` inside a function.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_GLOB`",
			input: `setopt NO_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1867",
					Message: "`setopt NO_GLOB` disables glob expansion — `rm *.log` chases the literal `*.log`, `for f in *.txt` loops once. Quote specific args or scope with `LOCAL_OPTIONS` inside a function.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1867")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
