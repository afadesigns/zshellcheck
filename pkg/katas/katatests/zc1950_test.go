// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1950(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `tune2fs -l $DEV` (read only)",
			input:    `tune2fs -l $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `tune2fs -m 1 $DEV` (tiny but non-zero reserve)",
			input:    `tune2fs -m 1 $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tune2fs -O ^has_journal $DEV`",
			input: `tune2fs -O ^has_journal $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1950",
					Message: "`tune2fs -O ^has_journal` strips the journal — crash recovery needs a full `fsck -y` and may truncate files. Keep the default.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `tune2fs -m 0 $DEV`",
			input: `tune2fs -m 0 $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1950",
					Message: "`tune2fs -m 0` zeroes the root reserve — a full fs leaves no headroom for `journald`/`apt`/root shells. Keep the default.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1950")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
