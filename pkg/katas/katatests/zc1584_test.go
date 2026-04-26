// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1584(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sudo -u root cmd",
			input:    `sudo -u root cmd`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — sudo env VAR=1 cmd",
			input:    `sudo env VAR=1 cmd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sudo -E cmd",
			input: `sudo -E cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1584",
					Message: "`sudo -E` carries the caller's PATH / LD_PRELOAD / … into the privileged process. Use `env_keep` in sudoers or explicit `sudo env VAR=… cmd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sudo -E -u svc cmd",
			input: `sudo -E -u svc cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1584",
					Message: "`sudo -E` carries the caller's PATH / LD_PRELOAD / … into the privileged process. Use `env_keep` in sudoers or explicit `sudo env VAR=… cmd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1584")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
