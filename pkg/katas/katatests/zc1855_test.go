// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1855(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `echo ${(k)groups}` (Zsh-native)",
			input:    `echo ${(k)groups}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `echo GROUPSIZE` (unrelated literal)",
			input:    `echo GROUPSIZE`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `echo $GROUPS`",
			input: `echo $GROUPS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1855",
					Message: "`$GROUPS` is a Bash-only array — Zsh populates `$groups` (associative name→GID) instead. Iterate `${(k)groups}` for names or `${(v)groups}` for GIDs, or fall back to `id -Gn`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `printf '%s\\n' \"${GROUPS[@]}\"`",
			input: `printf '%s\n' "${GROUPS[@]}"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1855",
					Message: "`$GROUPS` is a Bash-only array — Zsh populates `$groups` (associative name→GID) instead. Iterate `${(k)groups}` for names or `${(v)groups}` for GIDs, or fall back to `id -Gn`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1855")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
