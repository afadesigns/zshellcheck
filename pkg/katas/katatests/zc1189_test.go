// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1189(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid source file",
			input:    `source /etc/profile`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid source /dev/stdin",
			input: `source /dev/stdin`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1189",
					Message: "Avoid `source /dev/stdin`. Use `eval \"$(cmd)\"` for direct evaluation. `/dev/stdin` sourcing is fragile across platforms.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1189")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
