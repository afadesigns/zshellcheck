// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1954(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setfattr -n user.comment -v 'hello' /tmp/f`",
			input:    `setfattr -n user.comment -v 'hello' /tmp/f`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `getfattr -d /tmp/f` (read-only sibling)",
			input:    `getfattr -d /tmp/f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setfattr -n security.capability -v $VAL /usr/local/bin/app`",
			input: `setfattr -n security.capability -v $VAL /usr/local/bin/app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1954",
					Message: "`setfattr -n security.capability` writes the raw kernel xattr — bypasses `setcap`/`chcon`/`evmctl` validation and audit. Use the purpose-built tool.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setfattr -n security.selinux -v $CTX /etc/app`",
			input: `setfattr -n security.selinux -v $CTX /etc/app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1954",
					Message: "`setfattr -n security.selinux` writes the raw kernel xattr — bypasses `setcap`/`chcon`/`evmctl` validation and audit. Use the purpose-built tool.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1954")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
