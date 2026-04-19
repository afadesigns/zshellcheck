package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1825(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `scp src user@host:dst` (default SFTP on OpenSSH 9+)",
			input:    `scp src user@host:dst`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `scp -r dir user@host:/path`",
			input:    `scp -r dir user@host:/path`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `scp -O src user@host:dst`",
			input: `scp -O src user@host:dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1825",
					Message: "`scp -O` forces the legacy SCP wire protocol — the one exposed to filename-injection (CVE-2020-15778 class). Drop `-O` (default SFTP is safer), or use `sftp` / upgrade the remote server.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1825")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
