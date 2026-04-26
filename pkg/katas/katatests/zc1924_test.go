// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1924(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `virsh list --all` (read-only domain list)",
			input:    `virsh list --all`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `virt-top -d 5` (live view)",
			input:    `virt-top -d 5`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `virt-cat -d mydomain /etc/shadow`",
			input: `virt-cat -d mydomain /etc/shadow`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1924",
					Message: "`virt-cat` reads/writes the VM disk directly from the host — bypasses in-guest permissions, audit, and LUKS; a live VM risks corruption from double-mount. Snapshot first, work on the clone, prefer in-guest tooling.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `guestmount -d mydomain -i /mnt`",
			input: `guestmount -d mydomain -i /mnt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1924",
					Message: "`guestmount` reads/writes the VM disk directly from the host — bypasses in-guest permissions, audit, and LUKS; a live VM risks corruption from double-mount. Snapshot first, work on the clone, prefer in-guest tooling.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1924")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
