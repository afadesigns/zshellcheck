// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1786(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mount.cifs //h/s /mnt -o credentials=/etc/cifs-creds`",
			input:    `mount.cifs //h/s /mnt -o credentials=/etc/cifs-creds`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `mount.cifs //h/s /mnt -o guest`",
			input:    `mount.cifs //h/s /mnt -o guest`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mount.cifs //h/s /mnt -o username=u,password=hunter2`",
			input: `mount.cifs //h/s /mnt -o username=u,password=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1786",
					Message: "`mount.cifs ... password=…` leaks the SMB password into argv / `ps` / `/proc/PID/cmdline`. Use `credentials=/path/to/creds` (mode 0600) or `$PASSWD` env var instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `mount -t cifs //h/s /mnt -o user=u,password=hunter2`",
			input: `mount -t cifs //h/s /mnt -o user=u,password=hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1786",
					Message: "`mount ... password=…` leaks the SMB password into argv / `ps` / `/proc/PID/cmdline`. Use `credentials=/path/to/creds` (mode 0600) or `$PASSWD` env var instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1786")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
