// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1863(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt CASE_GLOB` (explicit default)",
			input:    `setopt CASE_GLOB`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt CASE_GLOB`",
			input: `unsetopt CASE_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1863",
					Message: "`unsetopt CASE_GLOB` flips every later glob to case-insensitive — `rm *.log` sweeps `APP.LOG`, dispatchers keyed on case collisions. Keep the option on; use `(#i)pattern` per-glob when you need case-folding.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_CASE_GLOB`",
			input: `setopt NO_CASE_GLOB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1863",
					Message: "`setopt NO_CASE_GLOB` flips every later glob to case-insensitive — `rm *.log` sweeps `APP.LOG`, dispatchers keyed on case collisions. Keep the option on; use `(#i)pattern` per-glob when you need case-folding.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1863")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
