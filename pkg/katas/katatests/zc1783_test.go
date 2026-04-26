// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1783(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `podman rmi myimage:old`",
			input:    `podman rmi myimage:old`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `podman system df` (read only)",
			input:    `podman system df`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `nerdctl system prune` (no -a, no --volumes)",
			input:    `nerdctl system prune`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `podman system reset --force`",
			input: `podman system reset --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1783",
					Message: "`podman system reset` wipes every container artifact on the host — images, volumes, networks, pods. Use narrower removals (`rmi`, `volume rm`, scoped `prune`) against the specific resource you intend to delete.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `nerdctl system prune -af --volumes`",
			input: `nerdctl system prune -af --volumes`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1783",
					Message: "`nerdctl system prune -a --volumes` wipes every container artifact on the host — images, volumes, networks, pods. Use narrower removals (`rmi`, `volume rm`, scoped `prune`) against the specific resource you intend to delete.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1783")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
