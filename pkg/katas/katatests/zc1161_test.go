// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1161(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sha256sum with file",
			input:    `sha256sum file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid sha256sum in pipeline",
			input: `sha256sum -`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1161",
					Message: "Consider `zmodload zsh/sha256` or `zmodload zsh/md5` for hash operations. Zsh modules avoid spawning external hashing processes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1161")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
