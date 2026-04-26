// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1062(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "egrep usage",
			input: `egrep 'foo|bar' file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1062",
					Message: "`egrep` is deprecated. Use `grep -E` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "grep -E usage",
			input:    `grep -E 'foo|bar' file`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1062")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
