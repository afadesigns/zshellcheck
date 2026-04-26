// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1737(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `wpa_passphrase MySSID` (passphrase via stdin)",
			input:    `wpa_passphrase MySSID`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `wpa_passphrase MySSID < /run/secrets/wifi`",
			input:    `wpa_passphrase MySSID < /run/secrets/wifi`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `wpa_passphrase MySSID hunter2`",
			input: `wpa_passphrase MySSID hunter2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1737",
					Message: "`wpa_passphrase SSID PASSWORD` puts the Wi-Fi passphrase in argv — visible in `ps`, `/proc`, history. Drop the PASSWORD argument and pipe it via stdin (`wpa_passphrase SSID < /run/secrets/wifi`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1737")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
