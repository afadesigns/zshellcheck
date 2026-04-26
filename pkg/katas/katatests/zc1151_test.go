// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1151(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid cat with file",
			input:    `cat file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat -A",
			input: `cat -A file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1151",
					Message: "Avoid `cat -A` for inspecting non-printable characters. Use `od -c` or `hexdump -C` for reliable cross-platform output.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1151")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
