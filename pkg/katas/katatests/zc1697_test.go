// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1697(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — cryptsetup open without --allow-discards",
			input:    `cryptsetup open $DISK data`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — cryptsetup luksClose",
			input:    `cryptsetup luksClose data`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — cryptsetup open --allow-discards",
			input: `cryptsetup open --allow-discards $DISK data`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1697",
					Message: "`cryptsetup --allow-discards` leaks free-sector layout to anyone with raw-device access — drop it if offline-disk inspection is in scope, or document the trade-off.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — cryptsetup luksOpen --allow-discards",
			input: `cryptsetup luksOpen --allow-discards $DISK data`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1697",
					Message: "`cryptsetup --allow-discards` leaks free-sector layout to anyone with raw-device access — drop it if offline-disk inspection is in scope, or document the trade-off.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1697")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
