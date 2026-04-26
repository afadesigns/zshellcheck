// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1971(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt GLOBAL_RCS` (keeps default on)",
			input:    `setopt GLOBAL_RCS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NO_GLOBAL_RCS` (restores default)",
			input:    `unsetopt NO_GLOBAL_RCS`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt GLOBAL_RCS`",
			input: `unsetopt GLOBAL_RCS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1971",
					Message: "`unsetopt GLOBAL_RCS` tells Zsh to skip `/etc/zprofile`, `/etc/zshrc`, `/etc/zlogin`, `/etc/zlogout` — corp `PATH`/audit/umask/proxy config silently dropped. Keep on; scope pristine setup with `emulate -LR zsh` or `env -i zsh -f`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_GLOBAL_RCS`",
			input: `setopt NO_GLOBAL_RCS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1971",
					Message: "`setopt NO_GLOBAL_RCS` tells Zsh to skip `/etc/zprofile`, `/etc/zshrc`, `/etc/zlogin`, `/etc/zlogout` — corp `PATH`/audit/umask/proxy config silently dropped. Keep on; scope pristine setup with `emulate -LR zsh` or `env -i zsh -f`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1971")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
