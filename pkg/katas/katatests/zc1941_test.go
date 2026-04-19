package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1941(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `restic init --password-file /etc/restic.pass`",
			input:    `restic init --password-file /etc/restic.pass`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `restic snapshots`",
			input:    `restic snapshots`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `restic init --insecure-no-password now`",
			input: `restic init --insecure-no-password now`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1941",
					Message: "`restic --insecure-no-password` creates an unencrypted repo — every operator with read access to the backend can reassemble the backed-up filesystem. Use `--password-file` / `--password-command` with a real passphrase.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `restic backup /data --insecure-no-password`",
			input: `restic backup /data --insecure-no-password`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1941",
					Message: "`restic --insecure-no-password` creates an unencrypted repo — every operator with read access to the backend can reassemble the backed-up filesystem. Use `--password-file` / `--password-command` with a real passphrase.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1941")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
