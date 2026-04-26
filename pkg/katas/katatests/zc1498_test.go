// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1498(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — mount /mnt/data /mnt/backup",
			input:    `mount /mnt/data /mnt/backup`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — mount -o ro,remount /",
			input:    `mount -o ro,remount /`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — mount -o remount,rw /",
			input: `mount -o remount,rw /`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1498",
					Message: "`mount -o remount,rw /` makes a read-only system path writable — use ostree / systemd-sysext or fix /etc/fstab.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — mount -o rw,remount /boot",
			input: `mount -o rw,remount /boot`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1498",
					Message: "`mount -o remount,rw /boot` makes a read-only system path writable — use ostree / systemd-sysext or fix /etc/fstab.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1498")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
