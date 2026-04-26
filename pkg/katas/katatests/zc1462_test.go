// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1462(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run without --ipc",
			input:    `docker run alpine`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker run --ipc=shareable",
			input:    `docker run --ipc=shareable alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --ipc=host (equals form)",
			input: `docker run --ipc=host alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1462",
					Message: "`--ipc=host` shares host shared memory and SysV IPC with the container — trivial data theft and side-channel vector. Use the default private IPC.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — podman run --ipc host (space form)",
			input: `podman run --ipc host alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1462",
					Message: "`--ipc=host` shares host shared memory and SysV IPC with the container — trivial data theft and side-channel vector. Use the default private IPC.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1462")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
