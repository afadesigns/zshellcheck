// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1892(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `install -m 0755 foo /usr/local/bin/foo`",
			input:    `install -m 0755 foo /usr/local/bin/foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `install -m 0644 foo.conf /etc/foo.conf`",
			input:    `install -m 0644 foo.conf /etc/foo.conf`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `install -m 4755 foo /usr/local/bin/foo` (setuid)",
			input: `install -m 4755 foo /usr/local/bin/foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1892",
					Message: "`install -m 4755` sets setuid/setgid bits at install time — every execution becomes a privesc vector. Use `0755` and grant narrow caps with `setcap` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `install -m 2755 foo /usr/local/bin/foo` (setgid)",
			input: `install -m 2755 foo /usr/local/bin/foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1892",
					Message: "`install -m 2755` sets setuid/setgid bits at install time — every execution becomes a privesc vector. Use `0755` and grant narrow caps with `setcap` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1892")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
