// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1475(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — setcap cap_net_bind_service",
			input:    `setcap cap_net_bind_service+ep /usr/bin/node`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — setcap cap_sys_admin+ep",
			input: `setcap cap_sys_admin+ep /usr/bin/foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1475",
					Message: "`setcap` granting dangerous capability `cap_sys_admin` makes the binary a privesc vector for any executing user. Scope narrower or use a dedicated service user.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — setcap cap_dac_override+ep",
			input: `setcap cap_dac_override+ep /usr/bin/filedump`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1475",
					Message: "`setcap` granting dangerous capability `cap_dac_override` makes the binary a privesc vector for any executing user. Scope narrower or use a dedicated service user.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — setcap cap_setuid+ep",
			input: `setcap 'cap_setuid=+ep' /usr/bin/maybebad`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1475",
					Message: "`setcap` granting dangerous capability `cap_setuid` makes the binary a privesc vector for any executing user. Scope narrower or use a dedicated service user.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1475")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
