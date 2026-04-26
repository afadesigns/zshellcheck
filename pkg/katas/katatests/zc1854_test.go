// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1854(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `yum-config-manager --add-repo https://…` (TLS)",
			input:    `yum-config-manager --add-repo https://mirror.example/app.repo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `zypper addrepo https://…` (TLS)",
			input:    `zypper addrepo https://mirror.example/app app`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `yum-config-manager --add-repo http://…`",
			input: `yum-config-manager --add-repo http://mirror.example/app.repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1854",
					Message: "`yum-config-manager --add-repo http://mirror.example/app.repo` registers a plaintext repo — on-path attacker can substitute packages and strip GPG-check directives. Use `https://` and pin `gpgkey=file://` in the `.repo`.",
					Line:    1,
					Column:  21,
				},
			},
		},
		{
			name:  "invalid — `dnf config-manager --add-repo http://…`",
			input: `dnf config-manager --add-repo http://mirror.example/app.repo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1854",
					Message: "`dnf config-manager --add-repo http://mirror.example/app.repo` registers a plaintext repo — on-path attacker can substitute packages and strip GPG-check directives. Use `https://` and pin `gpgkey=file://` in the `.repo`.",
					Line:    1,
					Column:  21,
				},
			},
		},
		{
			name:  "invalid — `zypper addrepo http://…`",
			input: `zypper addrepo http://mirror.example/app app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1854",
					Message: "`zypper addrepo http://mirror.example/app` registers a plaintext repo — on-path attacker can substitute packages and strip GPG-check directives. Use `https://` and pin `gpgkey=file://` in the `.repo`.",
					Line:    1,
					Column:  8,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1854")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
