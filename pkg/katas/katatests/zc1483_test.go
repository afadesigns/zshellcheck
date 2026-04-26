// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1483(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pip install in venv",
			input:    `pip install requests`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — pip install --user",
			input:    `pip install --user requests`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pip install --break-system-packages",
			input: `pip install --break-system-packages requests`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1483",
					Message: "`--break-system-packages` installs into distro-managed paths and collides with apt/dnf. Use a venv, pipx, or uv/poetry instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — pip3 install --break-system-packages",
			input: `pip3 install --break-system-packages foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1483",
					Message: "`--break-system-packages` installs into distro-managed paths and collides with apt/dnf. Use a venv, pipx, or uv/poetry instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1483")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
