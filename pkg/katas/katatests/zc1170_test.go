// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1170(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid pushd -q",
			input:    `pushd -q /tmp`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid pushd without -q",
			input: `pushd /tmp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1170",
					Message: "Use `pushd -q` to suppress directory stack output in scripts. Without `-q`, the stack is printed on every call.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid popd without -q",
			input: `popd /tmp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1170",
					Message: "Use `popd -q` to suppress directory stack output in scripts. Without `-q`, the stack is printed on every call.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1170")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
