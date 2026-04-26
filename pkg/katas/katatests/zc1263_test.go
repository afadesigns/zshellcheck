// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1263(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid apt-get",
			input:    `apt-get install curl`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid apt in script",
			input: `apt install curl`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1263",
					Message: "Use `apt-get` instead of `apt` in scripts. `apt` is for interactive use; `apt-get` has a stable scripting interface.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1263")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
