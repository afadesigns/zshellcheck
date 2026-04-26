// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1935(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `apt autoremove --dry-run` (preview)",
			input:    `apt autoremove --dry-run`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `apt autoremove` (no purge, config files kept)",
			input:    `apt autoremove`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `apt autoremove --purge -y`",
			input: `apt autoremove --purge -y`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1935",
					Message: "`apt autoremove` strips packages the resolver thinks are unused plus their configs — uproots packages installed manually but never `apt-mark manual`-ed. Dry-run first, mark keepers, drop `--purge` in CI.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `zypper rm --clean-deps foo`",
			input: `zypper rm --clean-deps foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1935",
					Message: "`zypper autoremove` strips packages the resolver thinks are unused plus their configs — uproots packages installed manually but never `apt-mark manual`-ed. Dry-run first, mark keepers, drop `--purge` in CI.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1935")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
