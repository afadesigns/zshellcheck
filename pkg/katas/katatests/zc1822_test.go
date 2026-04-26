// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1822(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `csrutil status` (read only)",
			input:    `csrutil status`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `spctl --status`",
			input:    `spctl --status`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `csrutil disable`",
			input: `csrutil disable`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1822",
					Message: "`csrutil disable` disables macOS SIP / Gatekeeper / kext-consent — every malware analyst's favorite persistence primitive. Re-enable (`csrutil enable` in recovery, `spctl --master-enable`) and keep the default policy on.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `spctl kext-consent disable`",
			input: `spctl kext-consent disable`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1822",
					Message: "`spctl kext-consent disable` disables macOS SIP / Gatekeeper / kext-consent — every malware analyst's favorite persistence primitive. Re-enable (`csrutil enable` in recovery, `spctl --master-enable`) and keep the default policy on.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1822")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
