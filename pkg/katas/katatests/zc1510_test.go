// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1510(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — auditctl -w /etc/passwd",
			input:    `auditctl -w /etc/passwd -p wa`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — auditctl -e 1",
			input:    `auditctl -e 1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — auditctl -e 0",
			input: `auditctl -e 0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1510",
					Message: "`auditctl -e 0` disables audit subsystem — anti-forensics tactic. Use `-e 2` for a reboot-locked maintenance window instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — auditctl -D",
			input: `auditctl -D`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1510",
					Message: "`auditctl -D` deletes every audit rule — anti-forensics tactic. Use `-e 2` for a reboot-locked maintenance window instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1510")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
