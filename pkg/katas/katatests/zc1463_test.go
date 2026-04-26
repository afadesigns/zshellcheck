// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1463(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run without --userns",
			input:    `docker run alpine`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker run --userns=keep-id",
			input:    `podman run --userns=keep-id alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --userns=host",
			input: `docker run --userns=host alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1463",
					Message: "`--userns=host` disables user-namespace remap — UID 0 in the container == UID 0 on the host. Leave the default remap on.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — podman run --userns host",
			input: `podman run --userns host alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1463",
					Message: "`--userns=host` disables user-namespace remap — UID 0 in the container == UID 0 on the host. Leave the default remap on.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1463")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
