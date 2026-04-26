// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1539(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — parted -s DISK print",
			input:    `parted -s $DISK print`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — parted DISK mklabel (interactive)",
			input:    `parted $DISK mklabel gpt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — parted -s DISK mklabel gpt",
			input: `parted -s $DISK mklabel gpt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1539",
					Message: "`parted -s <disk> mklabel` bypasses the confirmation prompt — a typo in the disk variable silently repartitions the wrong device.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — parted -s DISK rm 1",
			input: `parted -s $DISK rm 1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1539",
					Message: "`parted -s <disk> rm` bypasses the confirmation prompt — a typo in the disk variable silently repartitions the wrong device.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1539")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
