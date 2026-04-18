package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1678(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — borg init --encryption=repokey",
			input:    `borg init --encryption=repokey-blake2 /backup`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — borg list (different subcommand)",
			input:    `borg list /backup`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — borg init --encryption=none joined",
			input: `borg init --encryption=none /backup`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1678",
					Message: "`borg init --encryption=none` leaves archives unauthenticated and readable — use `--encryption=repokey-blake2` (or `keyfile-blake2`) and store the passphrase in `BORG_PASSPHRASE_FILE`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — borg init -e none",
			input: `borg init -e none /backup`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1678",
					Message: "`borg init --encryption=none` leaves archives unauthenticated and readable — use `--encryption=repokey-blake2` (or `keyfile-blake2`) and store the passphrase in `BORG_PASSPHRASE_FILE`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1678")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
