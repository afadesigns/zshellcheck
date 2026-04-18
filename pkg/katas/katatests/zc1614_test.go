package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1614(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — expect driving a non-auth dialog",
			input:    `expect -c 'spawn lftp host; expect lftp; send "ls\r"; interact'`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — no expect in use",
			input:    `ssh -i key host cmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — expect with password",
			input: `expect -c 'spawn ssh user@host; expect password; send "s3cret\r"; interact'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1614",
					Message: "`expect` script contains `password` / `passphrase` — the full argv lands in `ps` and audit logs. Switch to key-based auth, or read the credential from a protected file the expect script opens.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — expect with passphrase",
			input: `expect -c 'spawn ssh-keygen -p -f key; expect passphrase; send "x\r"'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1614",
					Message: "`expect` script contains `password` / `passphrase` — the full argv lands in `ps` and audit logs. Switch to key-based auth, or read the credential from a protected file the expect script opens.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1614")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
