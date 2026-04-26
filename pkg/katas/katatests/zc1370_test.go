// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1370(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — non-yes command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — yes str",
			input: `yes banana`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1370",
					Message: "Prefer Zsh `repeat N { print str }` over `yes str | head -n N` for producing N copies of a line. No external `yes` process, no pipe.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1370")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
