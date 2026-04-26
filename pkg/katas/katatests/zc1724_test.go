// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1724(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `pacman -Syu package` (full upgrade then install)",
			input:    `pacman -Syu package`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pacman -S package` (install without DB refresh)",
			input:    `pacman -S package`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `pacman -Sy` (refresh DB, no install)",
			input:    `pacman -Sy`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `pacman -Sy package`",
			input: `pacman -Sy package`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1724",
					Message: "`pacman -Sy <pkg>` is a partial-upgrade footgun — refresh the DB but install only one package against the newer metadata. Use `pacman -Syu` first, then `pacman -S <pkg>` (or `pacman -Syu --noconfirm <pkg>` to keep it atomic).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1724")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
