// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1441(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker ps",
			input:    `docker ps`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker prune (no -af)",
			input:    `docker system prune`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker system prune -af",
			input: `docker system prune -a -f`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1441",
					Message: "`docker prune -af` / `-a --force` deletes all unused resources without prompt. Scope with `--filter` or target one resource type.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — docker prune -af combined",
			input: `docker prune -af`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1441",
					Message: "`docker prune -af` / `-a --force` deletes all unused resources without prompt. Scope with `--filter` or target one resource type.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1441")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
