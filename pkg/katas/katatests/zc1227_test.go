// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1227(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid curl -fsSL",
			input:    `curl -fsSL https://example.com/file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid curl without URL",
			input:    `curl -s localhost`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid curl without -f",
			input: `curl -sL https://example.com/install.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1227",
					Message: "Use `curl -f` to fail on HTTP errors. Without `-f`, curl silently returns error pages (404, 500) as if they were successful.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1227")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
