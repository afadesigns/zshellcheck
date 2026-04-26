// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1340(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — non-shuf command",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — shuf -n 1",
			input: `shuf -n 1 file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1340",
					Message: "Avoid `shuf` for random selection — use Zsh `${array[RANDOM%$#array+1]}` with `$RANDOM` for in-shell randomness without spawning an external.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1340")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
