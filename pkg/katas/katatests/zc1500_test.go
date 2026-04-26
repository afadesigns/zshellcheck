// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1500(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — systemctl status sshd",
			input:    `systemctl status sshd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — systemctl edit --no-edit sshd",
			input:    `systemctl edit --no-edit sshd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — systemctl edit sshd",
			input: `systemctl edit sshd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1500",
					Message: "`systemctl edit` opens $EDITOR and waits for the user. Use a drop-in `/etc/systemd/system/<unit>.d/*.conf` + `daemon-reload` in scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — systemctl edit --full myapp.service",
			input: `systemctl edit --full myapp.service`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1500",
					Message: "`systemctl edit` opens $EDITOR and waits for the user. Use a drop-in `/etc/systemd/system/<unit>.d/*.conf` + `daemon-reload` in scripts.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1500")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
