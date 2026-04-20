package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1965(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `systemd-cryptenroll --wipe-slot=recovery $DEV`",
			input:    `systemd-cryptenroll --wipe-slot=recovery $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `systemd-cryptenroll --tpm2-device=auto $DEV`",
			input:    `systemd-cryptenroll --tpm2-device=auto $DEV`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `systemd-cryptenroll $DEV --wipe-slot=all`",
			input: `systemd-cryptenroll $DEV --wipe-slot=all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1965",
					Message: "`systemd-cryptenroll --wipe-slot=all` wipes every LUKS key slot (passphrase/recovery/TPM2/FIDO2) in one call. Enrol the new slot first, wipe a specific index, back up the header with `cryptsetup luksHeaderBackup`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1965")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
