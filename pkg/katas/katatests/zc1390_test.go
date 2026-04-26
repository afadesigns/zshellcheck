// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1390(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo $GROUPS (scalar)",
			input:    `echo $GROUPS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo ${GROUPS[@]}",
			input: `echo ${GROUPS[@]}`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1390",
					Message: "Zsh `$GROUPS` is a scalar (primary GID), not an array. For all group IDs use `${(k)groups}` (after `zmodload zsh/parameter`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1390")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
