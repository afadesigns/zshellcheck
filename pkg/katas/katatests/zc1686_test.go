// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1686(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — bare compinit",
			input:    `compinit -d $XDG_CACHE_HOME/zcompdump`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — different command",
			input:    `compaudit`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — compinit -C",
			input: `compinit -C`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1686",
					Message: "`compinit -C` (skip-security-check) loads `$fpath` files that are writable by others — any user on the host can inject shell code. Run `compaudit`, fix permissions, then `compinit` without the flag.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — compinit -u",
			input: `compinit -u`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1686",
					Message: "`compinit -u` (load-insecure-files) loads `$fpath` files that are writable by others — any user on the host can inject shell code. Run `compaudit`, fix permissions, then `compinit` without the flag.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1686")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
