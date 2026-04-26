// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1355(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — echo without -E",
			input:    `echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — echo -E raw",
			input: `echo -E "$line"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1355",
					Message: "Use `print -r` instead of `echo -E` for raw output. `-E` is a Bash-ism and ignored by POSIX echo.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1355")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
