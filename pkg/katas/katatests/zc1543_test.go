// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1543(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — go install pkg@v1.2.3",
			input:    `go install github.com/foo/bar@v1.2.3`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — cargo install --git url --rev sha",
			input:    `cargo install --git https://example.com/foo --rev abc123 foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — cargo install foo (crates.io pin via crate version)",
			input:    `cargo install foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — go install pkg@latest",
			input: `go install github.com/foo/bar@latest`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1543",
					Message: "`go install github.com/foo/bar@latest` is unpinned — HEAD-of-default can change between runs. Pin to a version tag or commit hash for reproducibility.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — cargo install --git url (no rev)",
			input: `cargo install --git https://example.com/foo foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1543",
					Message: "`cargo install --git (no --rev/--tag/--branch)` is unpinned — HEAD-of-default can change between runs. Pin to a version tag or commit hash for reproducibility.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1543")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
