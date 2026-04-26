// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1213(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid apt-get -y",
			input:    `apt-get -y install curl`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid apt-get update",
			input:    `apt-get update`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid apt-get install without -y",
			input: `apt-get install curl`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1213",
					Message: "Use `apt-get -y` in scripts. Without `-y`, apt-get prompts for confirmation which hangs in non-interactive execution.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1213")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
