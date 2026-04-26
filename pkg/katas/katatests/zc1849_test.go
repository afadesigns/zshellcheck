// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1849(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt ALL_EXPORT` (explicit default)",
			input:    `unsetopt ALL_EXPORT`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt ALL_EXPORT`",
			input: `setopt ALL_EXPORT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1849",
					Message: "`setopt ALL_EXPORT` marks every later assignment for export — secrets like `password=...` leak into every child's env. Drop it; use explicit `export`, or scope inside a `LOCAL_OPTIONS` function.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_ALL_EXPORT`",
			input: `unsetopt NO_ALL_EXPORT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1849",
					Message: "`unsetopt NO_ALL_EXPORT` marks every later assignment for export — secrets like `password=...` leak into every child's env. Drop it; use explicit `export`, or scope inside a `LOCAL_OPTIONS` function.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1849")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
