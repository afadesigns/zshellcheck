// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1607(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — safe.directory scoped to a path",
			input:    `git config --global safe.directory /workspace/repo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — unrelated git config",
			input:    `git config user.email "me@example.com"`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — git config safe.directory '*'",
			input: `git config --global safe.directory '*'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1607",
					Message: "`git config safe.directory '*'` trusts every directory — defeats CVE-2022-24765 protection. List specific paths, or fix the ownership mismatch.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — git -c safe.directory=* status",
			input: `git -c safe.directory=* status`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1607",
					Message: "`git config safe.directory '*'` trusts every directory — defeats CVE-2022-24765 protection. List specific paths, or fix the ownership mismatch.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1607")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
