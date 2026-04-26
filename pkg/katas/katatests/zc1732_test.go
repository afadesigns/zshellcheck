// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1732(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `flatpak install --user org.example.App`",
			input:    `flatpak install --user org.example.App`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `flatpak override --filesystem=~/Documents org.example.App`",
			input:    `flatpak override --filesystem=~/Documents org.example.App`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `flatpak override --filesystem=host org.example.App`",
			input: `flatpak override --filesystem=host org.example.App`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1732",
					Message: "`flatpak override --filesystem=host` removes the Flatpak sandbox — the app gets unrestricted host-filesystem access. Grant a specific subdirectory (e.g. `--filesystem=~/Documents:ro`) or use Filesystem portals.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `flatpak override --filesystem=home org.example.App`",
			input: `flatpak override --filesystem=home org.example.App`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1732",
					Message: "`flatpak override --filesystem=home` removes the Flatpak sandbox — the app gets unrestricted host-filesystem access. Grant a specific subdirectory (e.g. `--filesystem=~/Documents:ro`) or use Filesystem portals.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `flatpak run --filesystem=host org.example.App`",
			input: `flatpak run --filesystem=host org.example.App`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1732",
					Message: "`flatpak run --filesystem=host` removes the Flatpak sandbox — the app gets unrestricted host-filesystem access. Grant a specific subdirectory (e.g. `--filesystem=~/Documents:ro`) or use Filesystem portals.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1732")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
