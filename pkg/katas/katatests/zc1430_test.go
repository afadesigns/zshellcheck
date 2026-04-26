// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1430(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sched in-shell",
			input:    `sched +1:00 cmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — at now",
			input: `at now + 1 minute`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1430",
					Message: "Prefer Zsh `sched` (from `zsh/sched`) for in-shell scheduling instead of `at`/`batch`. No daemon dependency, runs in the current shell's environment.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1430")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
