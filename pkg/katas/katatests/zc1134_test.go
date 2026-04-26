// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1134(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sleep 5",
			input:    `sleep 5`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid sleep 30",
			input:    `sleep 30`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sleep 0.1",
			input: `sleep 0.1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1134",
					Message: "Avoid `sleep 0.1` in loops. Short sleep intervals suggest busy-waiting. Consider event-driven alternatives like `inotifywait` or `zle -F`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1134")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
