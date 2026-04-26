// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1116(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid tee with append",
			input:    `tee -a logfile`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid simple tee",
			input: `tee output.log`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1116",
					Message: "Use Zsh multios (`setopt multios`) instead of `tee`. With multios, `cmd > file1 > file2` writes to both files without spawning tee.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid tee with multiple files",
			input: `tee file1 file2`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1116",
					Message: "Use Zsh multios (`setopt multios`) instead of `tee`. With multios, `cmd > file1 > file2` writes to both files without spawning tee.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1116")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
