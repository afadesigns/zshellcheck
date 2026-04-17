package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1590(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sshpass -e (env var)",
			input:    `sshpass -e ssh user@host cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — sshpass -f FILE",
			input:    `sshpass -f /run/secrets/pw ssh user@host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sshpass -p 'secret' ssh ...",
			input: `sshpass -p 'secret' ssh user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1590",
					Message: "`sshpass -p` places the password in argv — visible in `ps` / `/proc/<pid>/cmdline`. Switch to key-based auth, or at least use `sshpass -f FILE` / `sshpass -e`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sshpass -psecret ssh ...",
			input: `sshpass -psecret ssh user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1590",
					Message: "`sshpass -p` places the password in argv — visible in `ps` / `/proc/<pid>/cmdline`. Switch to key-based auth, or at least use `sshpass -f FILE` / `sshpass -e`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1590")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
