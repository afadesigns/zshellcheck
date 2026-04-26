// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1183(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ls -la",
			input:    `ls -la`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ls -t",
			input: `ls -t`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1183",
					Message: "Use Zsh glob qualifiers `*(om[1])` for newest file or `*(Om[1])` for oldest instead of `ls -t`. Glob qualifiers avoid spawning external processes.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid ls -ltr",
			input: `ls -ltr`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1183",
					Message: "Use Zsh glob qualifiers `*(om[1])` for newest file or `*(Om[1])` for oldest instead of `ls -t`. Glob qualifiers avoid spawning external processes.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1183")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
