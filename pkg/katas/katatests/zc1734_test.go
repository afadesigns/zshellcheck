// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1734(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `useradd alice`",
			input:    `useradd alice`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cat /etc/passwd` (read-only)",
			input:    `cat /etc/passwd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `cp /tmp/passwd /etc/passwd`",
			input: `cp /tmp/passwd /etc/passwd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1734",
					Message: "`cp /etc/passwd` bypasses the lock that `vipw` / `vigr` / `useradd` use on the user-identity files. Edit through `vipw -e` / `vigr -e`, or mutate a single entry with `useradd` / `usermod` / `passwd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `tee /etc/shadow`",
			input: `tee /etc/shadow`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1734",
					Message: "`tee /etc/shadow` bypasses the lock that `vipw` / `vigr` / `useradd` use on the user-identity files. Edit through `vipw -e` / `vigr -e`, or mutate a single entry with `useradd` / `usermod` / `passwd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `echo entry >> /etc/group`",
			input: `echo entry >> /etc/group`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1734",
					Message: "`>> /etc/group` bypasses the lock that `vipw` / `vigr` / `useradd` use on the user-identity files. Edit through `vipw -e` / `vigr -e`, or mutate a single entry with `useradd` / `usermod` / `passwd`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1734")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
