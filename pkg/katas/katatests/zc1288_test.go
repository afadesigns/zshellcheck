// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1288(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid typeset usage",
			input:    `typeset -A mymap`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid local usage",
			input:    `local myvar=hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid declare usage",
			input: `declare -A mymap`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1288",
					Message: "Use `typeset` instead of `declare` in Zsh scripts. `typeset` is the native Zsh idiom.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid declare with -i flag",
			input: `declare -i count`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1288",
					Message: "Use `typeset` instead of `declare` in Zsh scripts. `typeset` is the native Zsh idiom.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1288")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
