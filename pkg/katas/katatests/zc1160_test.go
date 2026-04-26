// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1160(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid curl",
			input:    `curl -o file.tar.gz https://example.com/file.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid wget",
			input: `wget https://example.com/file.tar.gz`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1160",
					Message: "Prefer `curl` over `wget` for portability. `curl` is pre-installed on macOS and most Linux distributions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1160")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
