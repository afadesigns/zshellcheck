// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1692(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — kexec -l load only",
			input:    `kexec -l /boot/vmlinuz --initrd=/boot/initrd.img`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — kexec -u unload",
			input:    `kexec -u`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — kexec -e",
			input: `kexec -e`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1692",
					Message: "`kexec -e` jumps to a preloaded kernel without firmware reboot — wtmp / auditd see nothing. Use `systemctl kexec` or a real `systemctl reboot` to keep the audit trail.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1692")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
