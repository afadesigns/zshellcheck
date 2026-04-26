// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1256(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid other command",
			input:    `mkdir /tmp/dir`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid mkfifo without trap",
			input: `mkfifo /tmp/mypipe`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1256",
					Message: "Set up `trap 'rm -f pipe' EXIT` after `mkfifo`. Named pipes persist on the filesystem and need explicit cleanup.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1256")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
