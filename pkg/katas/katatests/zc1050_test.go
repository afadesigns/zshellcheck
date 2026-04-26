// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1050(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "safe glob loop",
			input:    `for f in *.txt; do echo $f; done`,
			expected: []katas.Violation{},
		},
		{
			name:     "arithmetic loop",
			input:    `for (( i=0; i<5; i++ )); do echo $i; done`,
			expected: []katas.Violation{},
		},
		{
			name:  "loop over ls output",
			input: `for f in $(ls); do echo $f; done`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1050",
					Message: "Avoid iterating over `ls` output. Use globs (e.g. `*.txt`) to handle filenames with spaces correctly.",
					Line:    1,
					Column:  10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1050")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
