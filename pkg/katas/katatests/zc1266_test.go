// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1266(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid nproc",
			input:    `nproc`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat /proc/cpuinfo",
			input: `cat /proc/cpuinfo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1266",
					Message: "Use `nproc` instead of parsing `/proc/cpuinfo` for CPU count. `nproc` is portable and available on Linux and macOS (via coreutils).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1266")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
