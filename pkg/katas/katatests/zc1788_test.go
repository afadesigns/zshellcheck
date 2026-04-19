package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1788(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ssh user@host` (default config)",
			input:    `ssh user@host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ssh -F ~/.ssh/config user@host`",
			input:    `ssh -F ~/.ssh/config user@host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ssh -F /tmp/ssh.conf user@host`",
			input: `ssh -F /tmp/ssh.conf user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1788",
					Message: "`ssh -F /tmp/ssh.conf` loads an alternate config from a mutable path — a tamper on that file can pin `ProxyCommand` to arbitrary code. Keep the config in `~/.ssh/config` (or a repo-owned path with `0600` perms).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `scp -F/var/tmp/conf src host:dst` (attached form)",
			input: `scp -F/var/tmp/conf src host:dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1788",
					Message: "`scp -F /var/tmp/conf` loads an alternate config from a mutable path — a tamper on that file can pin `ProxyCommand` to arbitrary code. Keep the config in `~/.ssh/config` (or a repo-owned path with `0600` perms).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1788")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
