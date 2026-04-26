// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1572(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run -e LOG_LEVEL=info",
			input:    `docker run -e LOG_LEVEL=info alpine`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker run -e PASSWORD (no value, inherits)",
			input:    `docker run -e PASSWORD alpine`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker run --env-file secrets alpine",
			input:    `docker run --env-file secrets alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run -e PASSWORD=hunter2",
			input: `docker run -e PASSWORD=hunter2 alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1572",
					Message: "`-e PASSWORD=<value>` writes the secret into `docker inspect` and `/proc/1/environ`. Use `--env-file` 0600 or `--secret`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — podman run -e API_KEY=abc123",
			input: `podman run -e API_KEY=abc123 alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1572",
					Message: "`-e API_KEY=<value>` writes the secret into `docker inspect` and `/proc/1/environ`. Use `--env-file` 0600 or `--secret`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1572")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
