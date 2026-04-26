// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1224(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid /proc/meminfo",
			input:    `cat /proc/meminfo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid free",
			input: `free -m`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1224",
					Message: "Avoid parsing `free` output — its format varies across versions. Read `/proc/meminfo` directly for reliable memory information.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1224")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
