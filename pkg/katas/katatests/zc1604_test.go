// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1604(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — source explicit file",
			input:    `source /etc/bashrc`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — source variable path (no glob)",
			input:    `source $HOME/dotfiles/common.sh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — source /etc/profile.d/*.sh",
			input: `source /etc/profile.d/*.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1604",
					Message: "`source /etc/profile.d/*.sh` loads every matched file. One attacker-writable match is arbitrary code execution. Use explicit filenames.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — . $HOME/dotfiles/*.sh",
			input: `. $HOME/dotfiles/*.sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1604",
					Message: "`. $HOME/dotfiles/*.sh` loads every matched file. One attacker-writable match is arbitrary code execution. Use explicit filenames.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1604")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
