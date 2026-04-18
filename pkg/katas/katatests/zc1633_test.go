package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1633(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — passphrase-file",
			input:    `gpg -d --passphrase-file /run/secrets/gpg file.gpg`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — passphrase-fd",
			input:    `gpg -d --passphrase-fd 0 file.gpg`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — gpg -d --passphrase 'secret'",
			input: `gpg -d --passphrase 'secret' file.gpg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1633",
					Message: "`gpg --passphrase` puts the passphrase in argv — visible via `ps`. Use `--passphrase-file`, `--passphrase-fd`, or `--pinentry-mode=loopback` with the value on stdin.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — gpg2 -d --passphrase $PW",
			input: `gpg2 -d --passphrase $PW file.gpg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1633",
					Message: "`gpg --passphrase` puts the passphrase in argv — visible via `ps`. Use `--passphrase-file`, `--passphrase-fd`, or `--pinentry-mode=loopback` with the value on stdin.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1633")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
