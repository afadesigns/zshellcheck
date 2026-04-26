// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1616(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — fsfreeze -u (unfreeze)",
			input:    `fsfreeze -u /mnt/backup`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — different command",
			input:    `mount /mnt/backup`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — fsfreeze -f /mnt/backup",
			input: `fsfreeze -f /mnt/backup`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1616",
					Message: "`fsfreeze -f` freezes the mountpoint — every write hangs until `fsfreeze -u` runs. Wrap the call in `trap 'fsfreeze -u PATH' EXIT` so the thaw fires even on failure.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — fsfreeze -f $ROOTFS",
			input: `fsfreeze -f $ROOTFS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1616",
					Message: "`fsfreeze -f` freezes the mountpoint — every write hangs until `fsfreeze -u` runs. Wrap the call in `trap 'fsfreeze -u PATH' EXIT` so the thaw fires even on failure.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1616")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
