// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1333(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid TIMEFMT usage",
			input:    `echo $TIMEFMT`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid TIMEFORMAT usage",
			input: `echo $TIMEFORMAT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1333",
					Message: "Avoid `$TIMEFORMAT` in Zsh — use `$TIMEFMT` instead. Format specifiers differ between Bash and Zsh.",
					Line:    1,
					Column:  6,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1333")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
