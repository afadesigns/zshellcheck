// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1555(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chmod on own file",
			input:    `chmod 600 /tmp/myfile`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — chown on /etc/nginx",
			input:    `chown root:root /etc/nginx/nginx.conf`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chmod 666 /etc/shadow",
			input: `chmod 666 /etc/shadow`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1555",
					Message: "`chmod ... /etc/shadow` races the distro-managed tool — use passwd/chage/visudo or a config-management drop-in.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chown root:root /etc/sudoers",
			input: `chown root:root /etc/sudoers`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1555",
					Message: "`chown ... /etc/sudoers` races the distro-managed tool — use passwd/chage/visudo or a config-management drop-in.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chgrp shadow /etc/gshadow",
			input: `chgrp shadow /etc/gshadow`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1555",
					Message: "`chgrp ... /etc/gshadow` races the distro-managed tool — use passwd/chage/visudo or a config-management drop-in.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1555")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
