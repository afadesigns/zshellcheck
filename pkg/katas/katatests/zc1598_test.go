// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1598(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chmod 600 on /dev/kvm",
			input:    `chmod 600 /dev/kvm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — chmod 666 on /dev/null (safe)",
			input:    `chmod 666 /dev/null`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — chmod 666 on regular file (not /dev/)",
			input:    `chmod 666 /tmp/log`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chmod 666 /dev/kvm",
			input: `chmod 666 /dev/kvm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1598",
					Message: "`chmod 666 /dev/kvm` makes a sensitive device node world-writable — direct kernel access for every local user. Keep restrictive perms (600 / 660) and grant access via udev rules.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chmod 777 /dev/mem",
			input: `chmod 777 /dev/mem`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1598",
					Message: "`chmod 777 /dev/mem` makes a sensitive device node world-writable — direct kernel access for every local user. Keep restrictive perms (600 / 660) and grant access via udev rules.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1598")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
