package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1823(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `keytool -list -keystore trust.jks`",
			input:    `keytool -list -keystore trust.jks`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `keytool -import -alias ca -file ca.pem -keystore trust.jks` (prompt)",
			input:    `keytool -import -alias ca -file ca.pem -keystore trust.jks`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `keytool -import -noprompt -alias ca -file ca.pem -keystore trust.jks`",
			input: `keytool -import -noprompt -alias ca -file ca.pem -keystore trust.jks`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1823",
					Message: "`keytool -import -noprompt` pins a cert to the Java trust store without a fingerprint check. Drop `-noprompt`, verify with `keytool -printcert -file CERT`, and store (alias, SHA-256) pairs in an audited inventory.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `keytool -importcert -noprompt -file ca.pem -keystore cacerts`",
			input: `keytool -importcert -noprompt -file ca.pem -keystore cacerts`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1823",
					Message: "`keytool -import -noprompt` pins a cert to the Java trust store without a fingerprint check. Drop `-noprompt`, verify with `keytool -printcert -file CERT`, and store (alias, SHA-256) pairs in an audited inventory.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1823")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
