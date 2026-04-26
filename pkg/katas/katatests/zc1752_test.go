// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1752(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `pvcreate $DISK` (prompts kept)",
			input:    `pvcreate $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `vgcreate my_vg $DISK`",
			input:    `vgcreate my_vg $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `pvcreate -ff $DISK`",
			input: `pvcreate -ff $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1752",
					Message: "`pvcreate -ff` skips the LVM confirmation — a wrong device gets its filesystem / RAID / LVM signatures wiped. Inspect with `wipefs -n` + `lsblk -f` first, drop the flag, re-run after checking the target.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `pvcreate $DISK --yes`",
			input: `pvcreate $DISK --yes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1752",
					Message: "`pvcreate --yes` skips the LVM confirmation — a wrong device gets its filesystem / RAID / LVM signatures wiped. Inspect with `wipefs -n` + `lsblk -f` first, drop the flag, re-run after checking the target.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `vgcreate -y my_vg $DISK`",
			input: `vgcreate -y my_vg $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1752",
					Message: "`vgcreate -y` skips the LVM confirmation — a wrong device gets its filesystem / RAID / LVM signatures wiped. Inspect with `wipefs -n` + `lsblk -f` first, drop the flag, re-run after checking the target.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1752")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
