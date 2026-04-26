// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1461(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run without --pid",
			input:    `docker run alpine`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker run --pid=container:other",
			input:    `docker run --pid=container:abc alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --pid=host (equals form)",
			input: `docker run --pid=host alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1461",
					Message: "`--pid=host` shares the host PID namespace — container can signal and inspect every host process. Avoid outside debug tools.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — podman run --pid host (space form)",
			input: `podman run --pid host alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1461",
					Message: "`--pid=host` shares the host PID namespace — container can signal and inspect every host process. Avoid outside debug tools.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1461")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
