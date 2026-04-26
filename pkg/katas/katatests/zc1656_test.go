// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1656(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — rsync with plain ssh",
			input:    `rsync -e ssh src host:dst`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — rsync without -e",
			input:    `rsync -a src dst`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — rsync -e "ssh -o StrictHostKeyChecking=no"`,
			input: `rsync -e "ssh -o StrictHostKeyChecking=no" src host:dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1656",
					Message: "`rsync -e 'ssh -o StrictHostKeyChecking=no'` disables host-key verification — MITM risk. Pre-provision `known_hosts` and keep `StrictHostKeyChecking=yes`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — rsync with UserKnownHostsFile=/dev/null`,
			input: `rsync -e "ssh -o UserKnownHostsFile=/dev/null" src host:dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1656",
					Message: "`rsync -e 'ssh -o StrictHostKeyChecking=no'` disables host-key verification — MITM risk. Pre-provision `known_hosts` and keep `StrictHostKeyChecking=yes`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1656")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
