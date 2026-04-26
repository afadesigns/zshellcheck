// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1128(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid touch with timestamp",
			input:    `touch -t 202301011200 file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid touch multiple files",
			input:    `touch file1 file2`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid touch single file",
			input: `touch newfile.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1128",
					Message: "Use `> file` instead of `touch file` to create an empty file. This avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1128")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
