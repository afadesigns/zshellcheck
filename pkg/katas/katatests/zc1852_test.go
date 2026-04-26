// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1852(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `firewall-cmd --panic-off foo`",
			input:    `firewall-cmd --panic-off foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `firewall-cmd --reload`",
			input:    `firewall-cmd --reload foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `firewall-cmd --panic-on >/dev/null` (mangled name)",
			input: `firewall-cmd --panic-on foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1852",
					Message: "`firewall-cmd --panic-on` drops every packet regardless of zone — an SSH-run call loses the session instantly. Use targeted zone rules; if you really need panic mode, gate behind `at now + N minutes … firewall-cmd --panic-off`.",
					Line:    1,
					Column:  15,
				},
			},
		},
		{
			name:  "invalid — `firewall-cmd \"\" --panic-on` (trailing flag)",
			input: `firewall-cmd "" --panic-on`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1852",
					Message: "`firewall-cmd --panic-on` drops every packet regardless of zone — an SSH-run call loses the session instantly. Use targeted zone rules; if you really need panic mode, gate behind `at now + N minutes … firewall-cmd --panic-off`.",
					Line:    1,
					Column:  18,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1852")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
