// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1358(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pwd without -P",
			input:    `pwd`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pwd -P",
			input: `pwd -P`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1358",
					Message: "Use `${PWD:P}` instead of `pwd -P` — the `P` modifier resolves symlinks and returns the canonical path without spawning an external.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1358")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
