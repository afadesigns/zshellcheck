// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1813(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `cryptsetup status cryptroot` (read only)",
			input:    `cryptsetup status cryptroot`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cryptsetup open $DEV cryptroot`",
			input:    `cryptsetup open $DEV cryptroot`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `cryptsetup luksFormat $DEV`",
			input: `cryptsetup luksFormat $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1813",
					Message: "`cryptsetup luksFormat` rewrites the LUKS header / device. Verify the target (`lsblk`), back up with `luksHeaderBackup`, and run on an unmounted volume with UPS.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `cryptsetup reencrypt $DEV`",
			input: `cryptsetup reencrypt $DEV`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1813",
					Message: "`cryptsetup reencrypt` rewrites the LUKS header / device. Verify the target (`lsblk`), back up with `luksHeaderBackup`, and run on an unmounted volume with UPS.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1813")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
