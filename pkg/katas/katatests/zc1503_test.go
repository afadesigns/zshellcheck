// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1503(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — groupadd mygroup",
			input:    `groupadd mygroup`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — groupadd -g 2000 mygroup",
			input:    `groupadd -g 2000 mygroup`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — groupadd -g 0 fakeroot",
			input: `groupadd -g 0 fakeroot`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1503",
					Message: "Creating a group with GID 0 duplicates the `root` group — hidden privesc. Pick an unused GID (see `getent group`) and scope via sudoers/polkit.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — groupmod -g0 service",
			input: `groupmod -g0 service`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1503",
					Message: "Creating a group with GID 0 duplicates the `root` group — hidden privesc. Pick an unused GID (see `getent group`) and scope via sudoers/polkit.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1503")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
