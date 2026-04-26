// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1680(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — vault file under /etc/ansible",
			input:    `ansible-playbook site.yml --vault-password-file=/etc/ansible/vault.pass`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — no vault file",
			input:    `ansible-playbook site.yml`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — vault file under /tmp joined",
			input: `ansible-playbook site.yml --vault-password-file=/tmp/vault.pass`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1680",
					Message: "`ansible-playbook --vault-password-file` under `/tmp/` / `/var/tmp/` / `/dev/shm/` — world-traversable, any local user can race-read it. Store the key mode-0400 under `/etc/ansible/` or supply via a `vault-password-client` helper.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — vault file under /dev/shm split",
			input: `ansible-playbook site.yml --vault-password-file /dev/shm/vault`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1680",
					Message: "`ansible-playbook --vault-password-file` under `/tmp/` / `/var/tmp/` / `/dev/shm/` — world-traversable, any local user can race-read it. Store the key mode-0400 under `/etc/ansible/` or supply via a `vault-password-client` helper.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1680")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
