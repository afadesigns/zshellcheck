// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1071(t *testing.T) {
	// t.Skip("Skipping ZC1071 tests due to parser limitation with array literals. See issue #41.")

	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid append self reference single",
			input: `arr=($arr)`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1071",
					Message: "Appending to an array using `arr=($arr ...)` is verbose and slower. Use `arr+=(...)` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "no self reference in array",
			input:    `arr=(a b c)`,
			expected: []katas.Violation{},
		},
		{
			name:     "simple assignment no array",
			input:    `x=5`,
			expected: []katas.Violation{},
		},
		{
			name:     "non-assignment operator",
			input:    `x + 5;`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1071")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
