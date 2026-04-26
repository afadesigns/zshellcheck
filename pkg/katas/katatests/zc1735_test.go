// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1735(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `efibootmgr -v` (inspect)",
			input:    `efibootmgr -v`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `efibootmgr -o 0001,0002` (reorder)",
			input:    `efibootmgr -o 0001,0002`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `efibootmgr -B`",
			input: `efibootmgr -B`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1735",
					Message: "`efibootmgr -B` deletes a UEFI boot entry — wrong BOOTNUM (or missing fallback) leaves the box at the UEFI shell on next reboot. Inspect `efibootmgr -v` first; demote via `-o NEW,ORDER` instead of deleting.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `efibootmgr -B -b 0001`",
			input: `efibootmgr -B -b 0001`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1735",
					Message: "`efibootmgr -B` deletes a UEFI boot entry — wrong BOOTNUM (or missing fallback) leaves the box at the UEFI shell on next reboot. Inspect `efibootmgr -v` first; demote via `-o NEW,ORDER` instead of deleting.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1735")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
