// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1632(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated command",
			input:    `rm /tmp/staging.log`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — shred -u FILE",
			input: `shred -u /tmp/secret.key`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1632",
					Message: "`shred` may not overwrite original blocks on ext4/btrfs/zfs. For guaranteed erasure, use full-disk encryption with key destruction, or `blkdiscard` when retiring an SSD.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — shred -n 3 file",
			input: `shred -n 3 /var/log/secret.log`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1632",
					Message: "`shred` may not overwrite original blocks on ext4/btrfs/zfs. For guaranteed erasure, use full-disk encryption with key destruction, or `blkdiscard` when retiring an SSD.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1632")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
