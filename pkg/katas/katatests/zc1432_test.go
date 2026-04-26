// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1432(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — passwd -l (lock)",
			input:    `passwd -l alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — passwd -d (delete)",
			input: `passwd -d alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1432",
					Message: "`passwd -d user` deletes the password — account becomes passwordless. Use `passwd -l user` to lock, or `usermod -L` + delete SSH keys to disable login.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1432")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
