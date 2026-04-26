// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1629(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — rsync with explicit path",
			input:    `rsync -a --rsync-path=/usr/bin/rsync src host:dst`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — rsync with no path override",
			input:    `rsync -a src host:dst`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — rsync --rsync-path="sudo rsync"`,
			input: `rsync -a --rsync-path="sudo rsync" src host:dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1629",
					Message: "`rsync --rsync-path='sudo rsync'` runs remote rsync under privilege escalation. Use a scoped sudoers rule on the remote host and keep the path explicit (`/usr/bin/rsync`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — rsync --rsync-path="doas rsync"`,
			input: `rsync -a --rsync-path="doas rsync" src host:dst`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1629",
					Message: "`rsync --rsync-path='doas rsync'` runs remote rsync under privilege escalation. Use a scoped sudoers rule on the remote host and keep the path explicit (`/usr/bin/rsync`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1629")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
