// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1764(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `git commit -m \"msg\"`",
			input:    `git commit -m "msg"`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `git commit -S -m \"msg\"` (signed)",
			input:    `git commit -S -m "msg"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `git commit --no-verify -m \"msg\"`",
			input: `git commit --no-verify -m "msg"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1764",
					Message: "`git commit --no-verify` skips pre-commit and commit-msg hooks — the last guardrail against secret leaks and broken tests. Fix the hook or carve a narrow exemption instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `git commit -n -m \"msg\"`",
			input: `git commit -n -m "msg"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1764",
					Message: "`git commit -n` skips pre-commit and commit-msg hooks — the last guardrail against secret leaks and broken tests. Fix the hook or carve a narrow exemption instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1764")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
