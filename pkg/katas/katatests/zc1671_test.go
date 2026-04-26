// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1671(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — install -m 755",
			input:    `install -m 755 src /usr/local/bin/dst`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — mkdir -m 0755",
			input:    `mkdir -m 0755 /opt/dir`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — install -m 777",
			input: `install -m 777 src /usr/local/bin/dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1671",
					Message: "`install -m 777` creates a world-writable target — drop the world-write bit (e.g. `0755` / `0644`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — mkdir -m 0666 (parser normalizes to 438)",
			input: `mkdir -m 0666 /shared`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1671",
					Message: "`mkdir -m 438` creates a world-writable target — drop the world-write bit (e.g. `0755` / `0644`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1671")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
