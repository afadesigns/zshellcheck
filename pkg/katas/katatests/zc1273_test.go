// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1273(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid grep -q usage",
			input:    `grep -q pattern file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid grep without redirect",
			input:    `grep pattern file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid grep redirected to /dev/null",
			input: `grep pattern file.txt /dev/null`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1273",
					Message: "Use `grep -q` instead of redirecting to `/dev/null`. It is faster and more idiomatic.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1273")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
