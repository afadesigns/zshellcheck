// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1521(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — strace -e trace=openat cmd",
			input:    `strace -e trace=openat ls`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — strace -f cmd",
			input: `strace -f ls`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1521",
					Message: "`strace` without `-e` captures every syscall including secrets in read/write buffers. Scope with `-e trace=<set>`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — strace cmd (bare)",
			input: `strace ls`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1521",
					Message: "`strace` without `-e` captures every syscall including secrets in read/write buffers. Scope with `-e trace=<set>`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1521")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
