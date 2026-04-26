// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1613(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh-keygen on host key",
			input:    `ssh-keygen -l -f /etc/ssh/ssh_host_rsa_key`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — cat public host key",
			input:    `cat /etc/ssh/ssh_host_rsa_key.pub`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — cat /etc/ssh/ssh_host_ed25519_key",
			input: `cat /etc/ssh/ssh_host_ed25519_key`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1613",
					Message: "Reading `/etc/ssh/ssh_host_ed25519_key` through a text tool copies private-key material into the process and often into logs / scrollback. Use `ssh-keygen -l -f` for metadata, or pass the path directly to the consumer.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — grep PRIVATE $HOME/.ssh/id_rsa",
			input: `grep PRIVATE $HOME/.ssh/id_rsa`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1613",
					Message: "Reading `$HOME/.ssh/id_rsa` through a text tool copies private-key material into the process and often into logs / scrollback. Use `ssh-keygen -l -f` for metadata, or pass the path directly to the consumer.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1613")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
