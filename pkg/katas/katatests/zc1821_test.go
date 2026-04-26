// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1821(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `diskutil list` (read only)",
			input:    `diskutil list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `diskutil info $DISK`",
			input:    `diskutil info $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `diskutil eraseDisk JHFS+ NewVol $DISK`",
			input: `diskutil eraseDisk JHFS+ NewVol $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1821",
					Message: "`diskutil eraseDisk` reformats the whole disk. Resolve the target by `diskutil info -plist` / mount-point (not by index), run `diskutil list` immediately before, and require a typed confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `diskutil secureErase 0 $DISK`",
			input: `diskutil secureErase 0 $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1821",
					Message: "`diskutil secureErase` overwrites every block, no undo. Resolve the target by `diskutil info -plist` / mount-point (not by index), run `diskutil list` immediately before, and require a typed confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1821")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
