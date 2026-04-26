// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1628(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — plain modprobe",
			input:    `modprobe nvme`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — lsmod",
			input:    `lsmod`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — insmod module.ko",
			input: `insmod evilmod.ko`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1628",
					Message: "`insmod` loads a kernel module bypassing depmod / blacklist — prefer `modprobe MODNAME` so system policy and signature checks apply.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — modprobe -f evilmod",
			input: `modprobe -f evilmod`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1628",
					Message: "`modprobe -f` ignores version-magic and kernel-mismatch checks — a mismatched module can crash or compromise the kernel. Drop the flag and fix the underlying version mismatch.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1628")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
