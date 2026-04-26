// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1249(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ssh-keygen -f",
			input:    `ssh-keygen -t ed25519 -f /tmp/key`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ssh-keygen without -f",
			input: `ssh-keygen -t rsa`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1249",
					Message: "Use `ssh-keygen -f /path/to/key -N ''` in scripts. Without `-f`, ssh-keygen prompts interactively for the file path.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1249")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
