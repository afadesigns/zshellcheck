// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1274(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid Zsh :t modifier",
			input:    `echo ${filepath:t}`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid non-basename command",
			input:    `dirname /path/to/file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid basename usage",
			input: `basename /path/to/file.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1274",
					Message: "Use Zsh parameter expansion `${var:t}` instead of `basename`. The `:t` modifier extracts the filename without forking a process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1274")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
