// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1857(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `cloud-init init` (boot-time init)",
			input:    `cloud-init init`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cloud-init status`",
			input:    `cloud-init status`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `cloud-init clean`",
			input: `cloud-init clean`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1857",
					Message: "`cloud-init clean` wipes `/var/lib/cloud/` boot state — the next reboot re-runs the user-data and overwrites operator changes (SSH host keys, hostname, `/etc/fstab`). Run interactively only.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `cloud-init clean --logs --reboot`",
			input: `cloud-init clean --logs --reboot`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1857",
					Message: "`cloud-init clean` wipes `/var/lib/cloud/` boot state — the next reboot re-runs the user-data and overwrites operator changes (SSH host keys, hostname, `/etc/fstab`). Run interactively only.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1857")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
