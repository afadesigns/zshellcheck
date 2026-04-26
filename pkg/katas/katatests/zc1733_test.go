// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1733(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `docker plugin install vieux/sshfs` (interactive prompt kept)",
			input:    `docker plugin install vieux/sshfs`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `docker plugin ls --grant-all-permissions` (not install)",
			input:    `docker plugin ls --grant-all-permissions`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `docker plugin install --grant-all-permissions vieux/sshfs`",
			input: `docker plugin install --grant-all-permissions vieux/sshfs`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1733",
					Message: "`docker plugin install --grant-all-permissions` accepts every capability the plugin requests — root-equivalent on the host. Walk the interactive prompt manually and pin the digest once vetted.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `docker plugin install vieux/sshfs --grant-all-permissions`",
			input: `docker plugin install vieux/sshfs --grant-all-permissions`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1733",
					Message: "`docker plugin install --grant-all-permissions` accepts every capability the plugin requests — root-equivalent on the host. Walk the interactive prompt manually and pin the digest once vetted.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1733")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
