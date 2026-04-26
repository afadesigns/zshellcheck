// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1456(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run -v local mount",
			input:    `docker run -v data:/app alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run -v /:/host",
			input: `docker run -v /:/host alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1456",
					Message: "`-v /:...` mounts the host root into the container — trivial container escape. Scope to specific paths.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1456")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
