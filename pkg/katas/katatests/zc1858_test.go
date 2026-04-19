package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1858(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ssh -c aes256-gcm@openssh.com host`",
			input:    `ssh -c aes256-gcm@openssh.com host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ssh host` (default ciphers)",
			input:    `ssh host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ssh -c 3des-cbc host`",
			input: `ssh -c 3des-cbc host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1858",
					Message: "`ssh ... 3des-cbc` forces a weak cipher with known plaintext-recovery / IV-reuse attacks. Leave cipher selection to OpenSSH defaults; if a legacy peer needs it, scope inside a `Host` block in `~/.ssh/config`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ssh -o Ciphers=arcfour,aes256-ctr host`",
			input: `ssh -o Ciphers=arcfour,aes256-ctr host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1858",
					Message: "`ssh ... arcfour` forces a weak cipher with known plaintext-recovery / IV-reuse attacks. Leave cipher selection to OpenSSH defaults; if a legacy peer needs it, scope inside a `Host` block in `~/.ssh/config`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1858")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
