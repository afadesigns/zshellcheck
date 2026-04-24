package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1898(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gpg --export KEYID`",
			input:    `gpg --export KEYID`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gpg --list-keys`",
			input:    `gpg --list-keys`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gpg --export-secret-keys KEYID` (leading)",
			input: `gpg --export-secret-keys KEYID`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1898",
					Message: "`gpg --export-secret-keys` writes the private key to stdout — one CI-log or wrong-tty redirect leaks it. Back up interactively on an air-gapped host, or write to a `umask 077` path and re-encrypt.",
					Line:    1,
					Column:  6,
				},
			},
		},
		{
			name:  "invalid — `gpg KEYID --export-secret-subkeys` (trailing)",
			input: `gpg KEYID --export-secret-subkeys`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1898",
					Message: "`gpg --export-secret-subkeys` writes the private key to stdout — one CI-log or wrong-tty redirect leaks it. Back up interactively on an air-gapped host, or write to a `umask 077` path and re-encrypt.",
					Line:    1,
					Column:  12,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1898")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
