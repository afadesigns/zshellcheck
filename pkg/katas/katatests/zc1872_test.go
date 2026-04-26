// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1872(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `badblocks -n $DISK`",
			input:    `badblocks -n $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `badblocks $DISK` (read-only)",
			input:    `badblocks $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `badblocks -w $DISK`",
			input: `badblocks -w $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1872",
					Message: "`badblocks -w` overwrites every sector of the target device — silent data wipe on a populated disk. Use `-n` (non-destructive) or gate destructive runs behind a confirmation and a fresh partition.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `badblocks -wsv $DISK` (combined)",
			input: `badblocks -wsv $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1872",
					Message: "`badblocks -w` overwrites every sector of the target device — silent data wipe on a populated disk. Use `-n` (non-destructive) or gate destructive runs behind a confirmation and a fresh partition.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1872")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
