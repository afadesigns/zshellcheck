// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1210(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid journalctl with --no-pager",
			input:    `journalctl --no-pager -u nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid journalctl without --no-pager",
			input: `journalctl -u nginx`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1210",
					Message: "Use `journalctl --no-pager` in scripts. Without it, journalctl invokes a pager that hangs in non-interactive execution.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1210")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
