// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1424(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — unrelated command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — mkfs.ext4",
			input: `mkfs.ext4 disk.img`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1424",
					Message: "`mkfs.ext4` formats / wipes a device — destroys data. Validate the target with `lsblk` / `blkid` first, and consider an interactive confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — mkswap",
			input: `mkswap swap.img`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1424",
					Message: "`mkswap` formats / wipes a device — destroys data. Validate the target with `lsblk` / `blkid` first, and consider an interactive confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1424")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
