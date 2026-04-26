// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1907(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sysctl -w fs.protected_symlinks=1` (re-enable)",
			input:    `sysctl -w fs.protected_symlinks=1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sysctl -w vm.swappiness=10` (unrelated)",
			input:    `sysctl -w vm.swappiness=10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `sysctl -w fs.protected_hardlinks=0`",
			input: `sysctl -w fs.protected_hardlinks=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1907",
					Message: "`sysctl -w fs.protected_hardlinks=0` re-enables hardlink following — classic /tmp-race escalation vector. Keep the default; scope any exception in a dedicated namespace.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `sysctl -w fs.suid_dumpable=2`",
			input: `sysctl -w fs.suid_dumpable=2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1907",
					Message: "`sysctl -w fs.suid_dumpable=2` re-enables SUID core-dump exposure (2 = root-readable) — classic /tmp-race escalation vector. Keep the default; scope any exception in a dedicated namespace.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1907")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
