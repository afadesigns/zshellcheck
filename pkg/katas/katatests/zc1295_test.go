// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1295(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid vared usage",
			input:    `vared myvar`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid read without -e",
			input:    `read -r myvar`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid read -e for editing",
			input: `read -e myvar`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1295",
					Message: "Use `vared` instead of `read -e` in Zsh. `vared` provides full ZLE editing support natively.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1295")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
