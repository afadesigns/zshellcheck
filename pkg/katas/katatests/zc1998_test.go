// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1998(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `tpm2 getcap algorithms`",
			input:    `tpm2 getcap algorithms`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `tpm2_pcrread sha256:0,1,2`",
			input:    `tpm2_pcrread sha256:0,1,2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tpm2_clear -c p`",
			input: `tpm2_clear -c p`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1998",
					Message: "`tpm2_clear` wipes the TPM storage hierarchy — every LUKS-TPM2 keyslot, `systemd-cryptenroll --tpm2-device` slot, and TPM-sealed TLS/sshd key is destroyed. No undo. Gate behind a recovery runbook.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `tpm2 clear -c p`",
			input: `tpm2 clear -c p`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1998",
					Message: "`tpm2 clear` wipes the TPM storage hierarchy — every LUKS-TPM2 keyslot, `systemd-cryptenroll --tpm2-device` slot, and TPM-sealed TLS/sshd key is destroyed. No undo. Gate behind a recovery runbook.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1998")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
