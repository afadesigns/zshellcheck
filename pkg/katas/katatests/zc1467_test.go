// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1467(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sysctl -w vm.swappiness=10",
			input:    `sysctl -w vm.swappiness=10`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — sysctl -w kernel.core_pattern=core",
			input:    `sysctl -w kernel.core_pattern=core`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — sysctl -w kernel.modprobe=/sbin/modprobe",
			input:    `sysctl -w kernel.modprobe=/sbin/modprobe`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sysctl -w 'kernel.core_pattern=|/tmp/x'",
			input: `sysctl -w 'kernel.core_pattern=|/tmp/x'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1467",
					Message: "Kernel hijack vector (kernel.core_pattern pipe handler) — next crash / module load runs attacker-supplied binary as root.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sysctl -w kernel.modprobe=/tmp/foo",
			input: `sysctl -w kernel.modprobe=/tmp/foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1467",
					Message: "Kernel hijack vector (kernel.modprobe override) — next crash / module load runs attacker-supplied binary as root.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1467")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
