// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1820(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `netplan try` (auto-reverting try)",
			input:    `netplan try`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `netplan get`",
			input:    `netplan get`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `netplan apply`",
			input: `netplan apply`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1820",
					Message: "`netplan apply` commits the YAML immediately — a mistake drops the admin SSH session with no automatic rollback. Run `netplan try` first (auto-reverts if no keypress within the timeout), then `netplan apply`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `netplan apply --debug`",
			input: `netplan apply --debug`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1820",
					Message: "`netplan apply` commits the YAML immediately — a mistake drops the admin SSH session with no automatic rollback. Run `netplan try` first (auto-reverts if no keypress within the timeout), then `netplan apply`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1820")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
