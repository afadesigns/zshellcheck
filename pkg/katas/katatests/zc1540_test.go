// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1540(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — cryptsetup luksOpen",
			input:    `cryptsetup luksOpen $DEV mapname`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — cryptsetup luksRemoveKey",
			input:    `cryptsetup luksRemoveKey $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — cryptsetup erase $DEV",
			input: `cryptsetup erase $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1540",
					Message: "`cryptsetup erase` wipes the LUKS header — ciphertext becomes unrecoverable. Back up the header first, or use luksRemoveKey/luksKillSlot for single-slot rotation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — cryptsetup luksErase $DEV",
			input: `cryptsetup luksErase $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1540",
					Message: "`cryptsetup luksErase` wipes the LUKS header — ciphertext becomes unrecoverable. Back up the header first, or use luksRemoveKey/luksKillSlot for single-slot rotation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1540")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
